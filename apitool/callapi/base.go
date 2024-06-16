// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package callapi provides simple way to call remote api server which is written
// with jsonapi package.
//
// For calling an API server written with pakcage jsonapi, [EP] and [NewEP] should
// solve your problem.
//
// For more complicated case, like multi-factor authetication, implementing your own
// [Endpoint] should solve your problem.
//
// Though not recommended, it's possible to call other APIs (Twitter, GCP, ...) by
// providing special designed [Encoder] and [Parser].
package callapi

// EP creates a Caller from http method and url, with common settings which is
// suitable to use with package jsonapi:
//
//   - If parameter is not nil, it is encoded by [json.Marshal] and content type is set to "application/json"
//   - request is sent by [http.DefaultClient]
//   - response is parsed by [DefaultParser]
//
// If your api server requires more configurations, like passing auth token or hmac
// signature with http header, you can:
//
//   - NewEP("POST", myApiUrl).With(aFunctionToSetHeader)
//   - Write your own Endpoint
func EP(method, url string) Caller {
	return NewEP(method, url).DefaultCaller()
}
