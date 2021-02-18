// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package apitool

import "github.com/raohwork/jsonapi"

// SessionData defines how you can access session data
//
// SessionData instance *MUST NOT* share across multiple goroutines.
type SessionData interface {
	// Returns session id
	ID() string
	// Delete a value
	Unset(key string)
	// Set a value
	Set(key string, val interface{})
	// Set a vallue that can be read once.
	SetOnce(key string, val interface{})
	// Get a value
	Get(key string) (val interface{}, ok bool)
	// Save session data to the store
	Save() error
	// Abandon this session
	Discard() error
}

// ErrSessionNotFound indicates there's no session registered with this request
type ErrSessionNotFound string

func (e ErrSessionNotFound) Error() string { return string(e) }

// SessionProvider defines how you can allocate session
//
// SessionProvider itself *MUST* be thread-safe, but SessionData loaded/created
// *MAY* be thread-unsafe.
type SessionProvider interface {
	// Returns existed session, or creates one if not exist.
	// Returns error if failed to get/create one due to internal error.
	Get(jsonapi.Request) (SessionData, error)
	// Requests to delete outdated session data.
	// It should be a NO-OP if the provider does not support it.
	GC()
}

// Session is a middleware to inject SessionData into request
func Session(p SessionProvider, key string) jsonapi.Middleware {
	return func(h jsonapi.Handler) (ret jsonapi.Handler) {
		return func(r jsonapi.Request) (data interface{}, err error) {
			sess, _ := p.Get(r)
			r = r.WithValue(key, sess)
			return h(r)
		}
	}
}

// GetSesion extracts SessionData which injected by middleware
func GetSession(r jsonapi.Request, key string) (sess SessionData, ok bool) {
	val := r.R().Context().Value(key)
	if val == nil {
		return
	}

	sess, ok = val.(SessionData)
	return
}
