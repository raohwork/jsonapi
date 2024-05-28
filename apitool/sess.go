// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package apitool

import (
	"context"

	"github.com/raohwork/jsonapi"
)

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

// SessionProvider defines how you can allocate session. It returns error only if
// the error is meant to report to user, otherwise it should returns empty data with
// no error.
//
// SessionProvider itself *MUST* be thread-safe, but SessionData loaded/created by
// it *MAY* be thread-unsafe.
type SessionProvider func(context.Context, jsonapi.Request) (SessionData, error)

// Session is a middleware to inject SessionData into context.
func Session(p SessionProvider, key string) jsonapi.Middleware {
	return func(h jsonapi.Handler) (ret jsonapi.Handler) {
		return func(ctx context.Context, r jsonapi.Request) (data interface{}, err error) {
			sess, err := p(ctx, r)
			if err != nil {
				return
			}

			ctx = context.WithValue(ctx, key, sess)
			return h(ctx, r)
		}
	}
}

// GetSesion extracts SessionData which injected by middleware
func GetSession(ctx context.Context, key string) (sess SessionData, ok bool) {
	val := ctx.Value(key)
	if val == nil {
		return
	}

	sess, ok = val.(SessionData)
	return
}
