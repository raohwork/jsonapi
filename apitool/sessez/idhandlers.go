// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package sessez

import (
	"net/http"
	"time"

	"github.com/raohwork/jsonapi"
)

// IDHandler defines how you pass session id to client, and how client pass it back
type IDHandler interface {
	// Returns session id if exist, or empty string if not found
	Get(jsonapi.Request) string
	// Pass session id to client, called by provider when creating new session.
	// Some handlers might have to ignore it.
	Set(r jsonapi.Request, id string) error
}

// InCookieSimple creates an IDHandler saves session id in http-only cookie
func InCookieSimple(key string) (ret IDHandler) {
	return inCookieSimple(key)
}

type inCookieSimple string

func (h inCookieSimple) Get(r jsonapi.Request) (id string) {
	req := r.R()
	c, err := req.Cookie(string(h))
	if err != nil {
		return
	}

	id = c.Value
	return
}

func (h inCookieSimple) Set(r jsonapi.Request, id string) (err error) {
	c := &http.Cookie{
		Name:     string(h),
		Value:    id,
		HttpOnly: true,
	}

	http.SetCookie(r.W(), c)
	return
}

// Encrypter defines how you encrypt session id
type Encrypter interface {
	Encrypt(string) string
	Decrypt(string) (string, error)
}

// InHeader passes session id in specified http header
type InHeader string

func (h InHeader) Get(r jsonapi.Request) (id string) {
	return r.R().Header.Get(string(h))
}

func (h InHeader) Set(r jsonapi.Request, id string) (err error) {
	r.W().Header().Set(string(h), id)
	return
}

// InCookie saves session id in cookie with more settings to enhance security
type InCookie struct {
	// Key name to store session id, required!
	Key string
	// Encrypt/decrypt session id if not nil
	Encrypter
	// Time-to-live for the cookie. The cookie is expired after TTL seconds.
	// 0 is http-only cookie, which will be destroyed when browser closed.
	TTL time.Duration
	// Refresh cookie expiration time each access.
	// It is ignored if TTL == 0.
	AutoRefresh bool
	// Enables secure cookie
	SecureFlag bool
}

func (h *InCookie) Get(r jsonapi.Request) (id string) {
	c, err := r.R().Cookie(h.Key)
	if err != nil {
		return
	}

	id = c.Value
	if h.Encrypter != nil {
		id, err = h.Decrypt(id)
		if err != nil {
			return
		}
	}

	if h.TTL == 0 || !h.AutoRefresh {
		return
	}

	c.Expires = time.Now().Add(h.TTL)
	http.SetCookie(r.W(), c)
	return
}

func (h *InCookie) Set(r jsonapi.Request, id string) (err error) {
	c := &http.Cookie{
		Name:   h.Key,
		Value:  id,
		Secure: h.SecureFlag,
	}

	if h.Encrypter != nil {
		c.Value = h.Encrypt(id)
	}

	if h.TTL == 0 {
		c.HttpOnly = true
	} else {
		c.Expires = time.Now().Add(h.TTL)
	}

	http.SetCookie(r.W(), c)
	return
}
