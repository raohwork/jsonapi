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

package sessez

import (
	"time"

	"github.com/raohwork/jsonapi"
	"github.com/raohwork/jsonapi/apitool"
)

// InternalData is the data structure used in storage
type InternalData struct {
	Val  interface{}
	Once bool
}

type sessionData struct {
	s    Store
	id   string
	data map[string]InternalData
	r    jsonapi.Request
	p    *sessionProvider
}

func (d *sessionData) ID() (id string) {
	return d.id
}

func (d *sessionData) Unset(key string) {
	delete(d.data, key)
}

func (d *sessionData) Set(key string, val interface{}) {
	d.data[key] = InternalData{Val: val}
}

func (d *sessionData) SetOnce(key string, val interface{}) {
	d.data[key] = InternalData{Val: val, Once: true}
}

func (d *sessionData) Get(key string) (ret interface{}, ok bool) {
	v, ok := d.data[key]
	if !ok {
		return
	}

	if v.Once {
		delete(d.data, key)
	}

	ret = v.Val
	return
}

func (d *sessionData) Save() (err error) {
	d.s.Set(d.id, d.data)
	if err = d.p.h.Set(d.r, d.id); err != nil {
		return
	}

	return
}

func (d *sessionData) Discard() (err error) {
	return d.s.Unset(d.id)
}

// New creates a SessionProvider
func New(saveID IDHandler, s Store, ttl time.Duration) (ret apitool.SessionProvider) {
	return &sessionProvider{
		saveID, s, ttl,
	}
}

type sessionProvider struct {
	h   IDHandler
	s   Store
	ttl time.Duration
}

func (p *sessionProvider) Get(r jsonapi.Request) (ret apitool.SessionData, err error) {
	id := p.h.Get(r)
	if id == "" {
		id, err = p.s.New(p.ttl)
		if err != nil {
			return
		}
	}

	i, err := p.s.Get(id)
	if err != nil {
		return
	}

	ret = &sessionData{
		s:    p.s,
		id:   id,
		data: i,
		r:    r,
		p:    p,
	}
	return
}

func (p *sessionProvider) GC() {
	p.s.GC()
}
