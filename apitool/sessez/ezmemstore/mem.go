// This file is part of jsonapi
//
// jsonapi is distributed in two licenses: The Mozilla Public License,
// v. 2.0 and the GNU Lesser Public License.
//
// jsonapi is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE.
//
// See LICENSE for further information.

// Package ezmemstore stores session data in memory
//
// You have to call GC periodically to release memory.
package ezmemstore

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/raohwork/jsonapi/apitool/sessez"
)

type memoryElement struct {
	*sync.Mutex // protects only data, lastUsed is not protected
	data        map[string]sessez.InternalData
	lastUsed    int64
	ttl         int64
}

func newMemEle(ttl int64) *memoryElement {
	return &memoryElement{
		&sync.Mutex{},
		map[string]sessez.InternalData{},
		time.Now().UnixNano(),
		ttl,
	}
}

func (e *memoryElement) isValid() bool {
	return time.Now().UnixNano() <= e.ttl+e.lastUsed
}

func (e *memoryElement) invalid() {
	// sould be now - ttl - 1, but set zero is faster
	e.lastUsed = 0
}

func (e *memoryElement) get() (data map[string]sessez.InternalData) {
	e.lastUsed = time.Now().UnixNano()
	return e.data
}

func (e *memoryElement) set(data map[string]sessez.InternalData) {
	e.lastUsed = time.Now().UnixNano()
	e.data = data
}

func generateRandomKey(size int, isExist func(string) bool) string {
	if m := size % 4; m != 0 {
		size += 4 - m
	}

	var ret string
	enc := base64.StdEncoding
	src := make([]byte, size/4*3)
	for ok := false; !ok; ok = isExist(ret) {
		// create a 64 bytes random data
		_, _ = rand.Read(src)

		ret = enc.EncodeToString(src)
	}

	return ret
}

type memoryStore struct {
	data   map[string]*memoryElement
	lock   *sync.Mutex // for allocate/release/gc
	gcing  bool
	lastgc int64
}

func (s *memoryStore) New(ttl time.Duration) (id string, err error) {
	id = generateRandomKey(32, func(id string) bool {
		s.lock.Lock()
		defer s.lock.Unlock()

		_, ok := s.data[id]
		if !ok {
			s.data[id] = newMemEle(int64(ttl))
		}

		return !ok
	})

	s.data[id] = newMemEle(int64(ttl))

	return
}

func (s *memoryStore) Unset(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.doDelete(id)
	return nil
}

func (s *memoryStore) doDelete(id string) {
	if data, ok := s.data[id]; ok {
		data.Lock()
		defer data.Unlock()
		data.invalid()

		delete(s.data, id)
	}
}

func (s *memoryStore) canGC() (ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	now := time.Now().UnixNano()
	if s.gcing {
		return
	}

	// at most 1 times in 1 sec
	if s.lastgc+int64(time.Second) >= now {
		return
	}

	s.gcing = true
	s.lastgc = now
	return true
}

// gc clears all invalid entries using Release, so no lock is required
func (s *memoryStore) GC() {
	if !s.canGC() {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for id, ele := range s.data {
		if ele.isValid() {
			continue
		}

		s.doDelete(id)
	}

	s.gcing = false
}

func (s *memoryStore) getElement(id string) (*memoryElement, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	e, ok := s.data[id]
	if !ok {
		return nil, errors.New("session not exists: " + id)
	}
	return e, nil
}

func (s *memoryStore) Get(id string) (ret map[string]sessez.InternalData, err error) {
	e, err := s.getElement(id)
	if err != nil {
		return
	}

	e.Lock()
	defer e.Unlock()
	if !e.isValid() {
		return nil, errors.New("session expired: " + id)
	}

	ret = e.get()
	return
}

func (s *memoryStore) Set(id string, data map[string]sessez.InternalData) error {
	e, err := s.getElement(id)
	if err != nil {
		return err
	}

	e.Lock()
	defer e.Unlock()
	if !e.isValid() {
		return errors.New("session expired: " + id)
	}
	e.set(data)

	return nil
}

// New creates a memory store
func New() sessez.Store {
	ret := &memoryStore{
		data: make(map[string]*memoryElement),
		lock: &sync.Mutex{},
	}
	return ret
}
