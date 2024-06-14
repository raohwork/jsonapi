// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package apitool

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/raohwork/jsonapi"
)

// EClient indicates something wrong at client side while calling remote jsonapi
type EClient struct {
	e error
}

func (e EClient) Error() string { return "api client error: " + e.e.Error() }
func (e EClient) Unwrap() error { return e.e }

// Client represnts a client binded to specified API endpoint. Zero value always
// returns error.
type Client struct {
	// endpoint. this is required
	URL string
	// http method to access this endpoint. this is required
	Method string
	// a optional function to modify request, like, set auth header
	Modifier func(*http.Request) *http.Request
	// http client to send request, nil = http.DefaultClient
	HTTPClient *http.Client
}

// Call sends param to the endpoint, and parses response with ParseResponse.
func (c *Client) Call(param, result any) error {
	return c.CallWithContext(context.TODO(), param, result)
}

// CallWithContext is same as Call, but the request is created with ctx.
func (c *Client) CallWithContext(ctx context.Context, param, result any) (err error) {
	resp, err := c.SendWithContext(ctx, param)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)
	return ParseResponse(resp, result)
}

// Send sends param to the endpoint. You have to parse response by your self.
func (c *Client) Send(ctx context.Context, param any) (resp *http.Response, err error) {
	return c.SendWithContext(context.TODO(), param)
}

// SendWithContext is same as Send, but the request is created with ctx.
func (c *Client) SendWithContext(ctx context.Context, param any) (resp *http.Response, err error) {
	cl := c.HTTPClient
	if cl == nil {
		cl = http.DefaultClient
	}

	buf, err := json.Marshal(param)
	if err != nil {
		return
	}
	body := bytes.NewReader(buf)
	req, err := http.NewRequestWithContext(ctx, c.Method, c.URL, body)
	if err != nil {
		return
	}
	if c.Modifier != nil {
		req = c.Modifier(req)
	}
	return cl.Do(req)
}

type callResp struct {
	Data   *json.RawMessage `json:"data"`
	Errors []jsonapi.ErrObj `json:"errors"`
}

// ParseResponse parses response of a jsonapi
//
// It's caller's response to close response body.
//
// If any io error or json decoding error occurred, an EClient is returned.
func ParseResponse(resp *http.Response, result interface{}) error {
	var res callResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return EClient{err}
	}

	if d := res.Data; d != nil {
		if err := json.Unmarshal([]byte(*d), result); err != nil {
			return EClient{err}
		}
	}

	if len(res.Errors) == 0 {
		return nil
	}

	return res.Errors[0].AsError()
}
