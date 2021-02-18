// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package apitool

import (
	"net/http"

	"github.com/raohwork/jsonapi"
)

// LogProvider defines what info a jsonapi logger can use
type LogProvider func(r *http.Request, data interface{}, err error)

// LogIn wraps handler, uses LogProvider p for logging purpose
//
//      jsonapi.With(
//          apitools.LogIn(apitool.JSONFormat(myLogger))
//      ).RegisterAll(myHandlerClass)
func LogIn(p LogProvider) jsonapi.Middleware {
	return jsonapi.Middleware(func(h jsonapi.Handler) jsonapi.Handler {
		return jsonapi.Handler(func(req jsonapi.Request) (interface{}, error) {
			data, err := h(req)
			p(req.R(), data, err)

			return data, err
		})
	})
}

// LogErrIn wraps handler, uses LogProvider p for logging purpose, but only for errors
//
//      jsonapi.With(
//          apitools.LogErrIn(apitool.JSONFormat(myLogger))
//      ).RegisterAll(myHandlerClass)
func LogErrIn(p LogProvider) jsonapi.Middleware {
	return jsonapi.Middleware(func(h jsonapi.Handler) jsonapi.Handler {
		return jsonapi.Handler(func(req jsonapi.Request) (interface{}, error) {
			data, err := h(req)
			if err != nil {
				p(req.R(), data, err)
			}

			return data, err
		})
	})
}
