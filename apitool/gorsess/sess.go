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

// Package gorsess wraps gorilla session to SessionProvider
package gorsess

import (
	"github.com/gorilla/sessions"
	"github.com/raohwork/jsonapi"
	"github.com/raohwork/jsonapi/apitool"
)

const index = "_flash_index"

func init() {
	// type guard
	_ = apitool.SessionData(&sessData{})
}

type sessData struct {
	s     *sessions.Session
	onceR map[string]interface{}
	onceW map[string]interface{}
	r     jsonapi.Request
}

func (s sessData) ID() string {
	return s.s.ID
}

func (s sessData) Unset(key string) {
	delete(s.s.Values, key)
}

func (s sessData) Set(key string, val interface{}) {
	s.s.Values[key] = val
}

func (s sessData) SetOnce(key string, val interface{}) {
	s.onceW[key] = val
}

func (s sessData) Get(key string) (val interface{}, ok bool) {
	val, ok = s.s.Values[key]
	if !ok {
		val, ok = s.onceR[key]
	}

	return
}

func (s sessData) Save() (err error) {
	keys := make([]string, 0, len(s.onceW))
	for k, v := range s.onceW {
		s.s.AddFlash(v, k)
		keys = append(keys, k)
	}

	s.s.AddFlash(keys, index)

	return s.s.Save(s.r.R(), s.r.W())
}

func (s sessData) Discard() (err error) {
	s.s.Values = map[interface{}]interface{}{}
	s.onceW = map[string]interface{}{}
	s.s.Options.MaxAge = -1
	return s.Save()
}

// New creates a SessionProvider using gorilla/sessions.Store
//
// It does not support on-demand garbage collecting (SessionProvider.GC()).
func New(store sessions.Store, name string) (ret apitool.SessionProvider) {
	return &provider{
		Store: store,
		Name:  name,
	}
}

type provider struct {
	Store sessions.Store
	Name  string
}

func (p *provider) Get(r jsonapi.Request) (ret apitool.SessionData, err error) {
	s, err := p.Store.Get(r.R(), p.Name)
	if err != nil {
		return
	}

	x := &sessData{
		s:     s,
		onceR: map[string]interface{}{},
		onceW: map[string]interface{}{},
		r:     r,
	}
	ret = x

	arr := s.Flashes(index)
	if len(arr) == 0 {
		return
	}

	keys, ok := arr[0].([]string)
	if !ok {
		return
	}

	vals := s.Flashes(keys...)
	for idx, k := range keys {
		x.onceR[k] = vals[idx]
	}
	return
}

func (p *provider) GC() {}
