// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package jsonapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Handler is easy to use entry for API developer.
//
// Just return something, and it will be encoded to JSON format and send to client.
// Or return an Error to specify http status code and error string.
//
//	func myHandler(ctx context.Context, req *Request) (interface{}, error) {
//	    var param paramType
//	    if err := req.Decode(&param); err != nil {
//	        return nil, jsonapi.E400.SetData("You must send parameters in JSON format.")
//	    }
//	    return doSomething(param), nil
//	}
//
// The context is derived from [http.Request.Context]. You should check if it is
// canceled before doing time/resource consuming job in general. There's a helper
// [IsCanceled].
//
// To redirect clients, return 301~303 status code and set Data property
//
//	return nil, jsonapi.E301.SetData("http://google.com")
//
// Redirecting depends on http.Redirect(). The data returned from handler will never
// write to ResponseWriter.
//
// This basically obey the http://jsonapi.org rules:
//
//   - Return {"data": your_data} if error == nil
//   - Return {"errors": [{"code": "error-code", "detail": "message"}]} if error returned
type Handler func(ctx context.Context, r Request) (interface{}, error)

// ServeHTTP implements net/http.Handler
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer io.Copy(io.Discard, r.Body)

	w.Header().Set("Content-Type", "application/json")
	res, err := h(r.Context(), FromHTTP(w, r))
	if IsCanceled(r.Context()) != nil {
		// connection is closed, do not output anything
		return
	}

	enc := json.NewEncoder(w)
	if err == nil {
		err = enc.Encode(map[string]interface{}{
			"data": res,
		})
		if err == nil {
			return
		}
		err = errors.New("Failed to marshal data")
	}

	resp := make(map[string]interface{})
	code := http.StatusInternalServerError
	if httperr, ok := err.(Error); ok {
		if httperr.EqualTo(ASIS) {
			if res != nil {
				switch x := res.(type) {
				case string:
					w.Write([]byte(x))
				case []byte:
					w.Write(x)
				case fmt.Stringer:
					w.Write([]byte(x.String()))
				default:
					fmt.Fprintf(w, "%v", res)
				}
			}
			return
		}
		code = httperr.Code
		if code >= 301 && code <= 303 && httperr.location != "" {
			// 301~303 redirect
			http.Redirect(w, r, httperr.location, code)
			return
		}

		w.WriteHeader(code)
		resp["errors"] = []*ErrObj{fromError(&httperr)}
		enc.Encode(resp)
		return
	}

	w.WriteHeader(code)
	resp["errors"] = []*ErrObj{{Detail: err.Error()}}
	enc.Encode(resp)
}
