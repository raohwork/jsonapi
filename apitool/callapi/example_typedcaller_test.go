// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

import (
	"context"
	"net/http"
)

type SpecialImplParam struct{}
type SpecialImplResp string

type BadImplParam struct{}
type BadImplResp string

type GoodImplParam struct{}
type GoodImplResp string

func authWith(token, secret string) func(*http.Request) (*http.Request, error) {
	return func(r *http.Request) (*http.Request, error) {
		r.Header.Set("MY-API-TOKEN", token)
		r.Header.Set("MY-API-SECRET", secret)
		return r, nil
	}
}

func NewClient(host, token, secret string) *MyAPIClient {
	return &MyAPIClient{
		b: Builder{
			Maker: func(method, path string) Endpoint {
				return DefaultEncoder().
					EP(method, host+path).
					With(authWith(token, secret))
			},
		},
	}
}

type MyAPIClient struct {
	b Builder
}

// see how bad it can be without using [Builder] for repeative code.
func (c *MyAPIClient) SpecialImpl(ctx context.Context, param SpecialImplParam) (*SpecialImplResp, error) {
	// host, token, secret should be stored in MyAPIClient if not using Builder
	host, token, secret := "", "", ""

	var ret SpecialImplResp
	err := NewEP(http.MethodPost, host+"/spec/impl").
		With(authWith(token, secret)).
		DefaultCaller().
		Call(ctx, param, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// imagine writing this code 20 times.....
func (c *MyAPIClient) BadImpl(ctx context.Context, param BadImplParam) (*BadImplResp, error) {
	var ret BadImplResp
	err := c.b.EP(http.MethodPost, "/bad/impl").Call(ctx, param, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// saves few key strokes
func (c *MyAPIClient) GoodImpl(ctx context.Context, param GoodImplParam) (*GoodImplResp, error) {
	return Typed[GoodImplParam, GoodImplResp](
		c.b.EP(http.MethodPost, "/bad/impl"),
	).Call(ctx, param)

	// if this method uses different set of request headers
	// return Use[GoodImplParam, GoodImplResp](
	// 	c.b.UseMaker(anotherMaker).EP(http.MethodPost, "/bad/impl"),
	// ).Call(ctx, param)
}

func ExampleTypedCaller() {
	// see methods of MyAPIClient for detail
	_ = NewClient("", "", "")
}
