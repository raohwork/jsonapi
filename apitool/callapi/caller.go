// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

import "context"

// Caller is a function calls to specific API endpoint.
//
// In general, you build an Endpoint (using Encoder or implement by yourself) and
// create Caller by applying Sender and Parser to it:
//
//	var result MyAPIResultType
//	param := MyAPIParam{ ... }
//	err := endpoint.SendBy(sneder).ParseWith(parser).Call(ctx, param, result)
type Caller func(ctx context.Context, param, result any) error

// Call calls the API with Caller c. It's identical to c(ctx, param, result).
func (c Caller) Call(ctx context.Context, param, result any) error {
	return c(ctx, param, result)
}

// TypedCaller is type-safe [Caller]. Note that returned type is pointer.
//
// It is designed to save some key strokes when writing API client method which
// supposed to return pointer type. For those methods returning value type:
//
//	func (c *MyClient) MyMethod(ctx context.Context, param MyParam) (ret MyResult, err error) {
//		err = c.builder.EP("POST", c.host+"/my/method").Call(ctx, param, &ret)
//		return
//	}
type TypedCaller[I, O any] func(context.Context, I) (*O, error)

// Call calls the API with TypedCaller c. It's identical to c(ctx, param).
func (c TypedCaller[I, O]) Call(ctx context.Context, param I) (*O, error) {
	return c(ctx, param)
}

// Typed creates [TypedCaller] from [Caller].
func Typed[I, O any](c Caller) TypedCaller[I, O] {
	return func(ctx context.Context, param I) (*O, error) {
		var ret O
		if err := c(ctx, param, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}
