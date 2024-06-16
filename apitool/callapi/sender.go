// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

import (
	"context"
	"net/http"
)

// Sender is s function to send request to server.
type Sender func(ctx context.Context, param any) (resp *http.Response, err error)

// ParseWith creates a Caller by parsing response returned by Sender s.
func (s Sender) ParseWith(p Parser) Caller {
	return func(ctx context.Context, param, result any) (err error) {
		resp, err := s(ctx, param)
		if err != nil {
			return
		}

		return p(resp, result)
	}
}

// Caller is a function calls to specific API endpoint.
//
// In general, you build an Endpoint (using Encoder or implement by yourself) and
// create Caller by applying Sender and Parser to it:
//
//	var result MyAPIResultType
//	param := MyAPIParam{ ... }
//	err := endpoint.By(sneder).With(parser).Call(ctx, param, result)
type Caller func(ctx context.Context, param, result any) error

// Call calls the API with Caller c. It's identical to c(ctx, param, result).
func (c Caller) Call(ctx context.Context, param, result any) error {
	return c(ctx, param, result)
}
