// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// DefaultEncoder returns default encoder, which is simple json.Marshal.
func DefaultEncoder() Encoder { return json.Marshal }

// Encoder is a function which can encodes param into specific format.
type Encoder func(v any) ([]byte, error)

// EP creates an Endpoint which uses the Encoder to build request body.
func (e Encoder) EP(method, url string) Endpoint {
	return func(ctx context.Context, param any) (req *http.Request, err error) {
		var body io.Reader
		if param != nil {
			buf, err := e(param)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(buf)
		}

		req, err = http.NewRequestWithContext(
			ctx, method, url, body,
		)
		if err != nil {
			return
		}

		if param != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		return
	}
}

// SlowSortEncoder is an Encoder which produces sorted json in super slow way.
//
// It's not recommended to use in production.
func SlowSortEncoder() Encoder { return slowSort }

func slowSort(v any) ([]byte, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var m interface{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}
