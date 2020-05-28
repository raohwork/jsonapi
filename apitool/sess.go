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

package apitool

import "github.com/raohwork/jsonapi"

// SessionData defines how you can access session data
type SessionData interface {
	// Returns session id
	ID() string
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
type SessionProvider interface {
	// Returns existed session, or creates one if not exist.
	// Returns error if failed to get/create one due to internal error.
	Get(jsonapi.Request) (SessionData, error)
	// Requests to delete outdated session data.
	// It should be a NO-OP if the provider does not support it.
	GC()
}

// Session is a middleware to inject session data into request
func Session(p SessionProvider, key string) jsonapi.Middleware {
	return func(h jsonapi.Handler) (ret jsonapi.Handler) {
		return func(r jsonapi.Request) (data interface{}, err error) {
			sess, _ := p.Get(r)
			r = r.WithValue(key, sess)
			return h(r)
		}
	}
}
