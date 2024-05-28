// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package apitest provides few tools helping you write tests
package apitest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/raohwork/jsonapi"
)

// NewRequest wraps httptest.NewRequest, use your data (encoded to JSON) as request body
//
// It also sets "Content-Type" to "application/json".
func NewRequest(method, target string, data interface{}) *http.Request {
	buf, _ := json.Marshal(data)
	return httptest.NewRequest(method, target, bytes.NewReader(buf))
}

// Modify creates a middleware that do some magic before running handler
func Modify(f func(jsonapi.Request) jsonapi.Request) jsonapi.Middleware {
	return func(h jsonapi.Handler) jsonapi.Handler {
		return func(ctx context.Context, r jsonapi.Request) (interface{}, error) {
			return h(ctx, f(r))
		}
	}
}

// Monitor creates a middleware that do some magic after running handler
func Monitor(f func(jsonapi.Request, interface{}, error)) jsonapi.Middleware {
	return func(h jsonapi.Handler) jsonapi.Handler {
		return func(ctx context.Context, r jsonapi.Request) (interface{}, error) {
			data, err := h(ctx, r)
			f(r, data, err)
			return data, err
		}
	}
}

// Test wraps your handler for test purpose
type Test jsonapi.Handler

// With creates new Test instance by wrapping the handler with the middleware
//
// It executes in REVERSE ORDER:
//
//	// order: m2 > m1 > h > m1 > m2
//	Test(h).With(m1).With(m2).Use(data)
func (t Test) With(m jsonapi.Middleware) Test {
	return Test(m(jsonapi.Handler(t)))
}

// UseRequest executes handler with specified request
func (t Test) UseRequest(req *http.Request) (interface{}, error) {
	defer req.Body.Close()

	w := httptest.NewRecorder()
	return t(req.Context(), jsonapi.FromHTTP(w, req))
}

// Use executes handler with your data
//
// The request address will be "/" and using POST method.
func (t Test) Use(data interface{}) (interface{}, error) {
	return t.UseRequest(
		NewRequest("POST", "/", data),
	)
}
