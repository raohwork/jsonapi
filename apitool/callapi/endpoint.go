// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

import (
	"context"
	"net/http"
)

// Endpoint represents the spec of an API endpoint, a function that creates http
// request which fulfills all requirements the endpoint need.
type Endpoint func(context.Context, any) (*http.Request, error)

// With wraps ep by modifying the request with f.
func (ep Endpoint) With(f func(*http.Request) (*http.Request, error)) Endpoint {
	return func(ctx context.Context, param any) (req *http.Request, err error) {
		req, err = ep(ctx, param)
		if err != nil {
			return
		}
		return f(req)
	}
}

// SendBy creates a Sender which creates request using ep and send it by cl.
func (ep Endpoint) SendBy(cl *http.Client) Sender {
	if cl == nil {
		cl = http.DefaultClient
	}

	return func(ctx context.Context, param any) (resp *http.Response, err error) {
		req, err := ep(ctx, param)
		if err != nil {
			return
		}
		return cl.Do(req)
	}
}

// DefaultCaller is shortcut to ep.SendBy(nil).ParseWith(DefaultParser)
func (ep Endpoint) DefaultCaller() Caller {
	return ep.SendBy(nil).ParseWith(DefaultParser)
}

// NewEP creates an Endpoint with default settings, see [EP] for detail.
func NewEP(method, url string) Endpoint {
	return DefaultEncoder().EP(method, url)
}
