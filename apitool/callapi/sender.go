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

// ParseWith creates a Caller by parsing response returned by Sender s. It uses
// [DefaultParser] if p == nil.
func (s Sender) ParseWith(p Parser) Caller {
	return func(ctx context.Context, param, result any) (err error) {
		resp, err := s(ctx, param)
		if err != nil {
			return
		}

		if p == nil {
			p = DefaultParser
		}
		return p(resp, result)
	}
}
