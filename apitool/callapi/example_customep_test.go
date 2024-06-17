// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
)

// MyEP creates an Endpoint which signs the request with hmac-sha1.
func MyEP(method, url string) Endpoint {
	return func(ctx context.Context, param any) (req *http.Request, err error) {
		buf, err := SlowSortEncoder()(param)
		if err != nil {
			return
		}

		req, err = http.NewRequestWithContext(
			ctx, method, url, bytes.NewReader(buf),
		)
		if err != nil {
			return
		}

		signer := hmac.New(sha1.New, []byte("my secret key"))
		io.WriteString(signer, req.URL.Path)
		io.WriteString(signer, ",my-client-id,")
		signer.Write(buf)
		sig := signer.Sum(nil)
		req.Header.Set("MY-SIGNATURE", hex.EncodeToString(sig))
		return
	}
}

func Auth(req *http.Request) (*http.Request, error) {
	req.Header.Set("MY-AUTH-TOKEN", "my secret token")
	return req, nil
}

func Example_customEndpoint() {
	// simple endpoint with auth token in header
	ep := NewEP(http.MethodGet, "https://example.com/api/my_endpoint").
		With(Auth).DefaultCaller()
	// var result MyResult
	// err = ep.Call(param, &result)

	// customized endpoint, with additional hmac signature in header
	ep = MyEP(http.MethodPost, "https://example.com/api/my_endpoint").
		With(Auth).DefaultCaller()
	// var result MyResult
	// err = ep.Call(param, &result)

	_ = ep
}
