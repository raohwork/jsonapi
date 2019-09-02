// This file is part of jsonapi
//
// jsonapi is distributed in two licenses: The Mozilla Public License,
// v. 2.0 and the GNU Lesser Public License.
//
// jsonapi is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE.
//
// See LICENSE for further information.

package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler is easy to use entry for API developer.
//
// Just return something, and it will be encoded to JSON format and send to client.
// Or return an Error to specify http status code and error string.
//
//     func myHandler(dec *json.Decoder, httpData *HTTP) (interface{}, error) {
//         var param paramType
//         if err := dec.Decode(&param); err != nil {
//             return nil, jsonapi.E400.SetData("You must send parameters in JSON format.")
//         }
//         return doSomething(param), nil
//     }
//
// To redirect clients, return 301~303 status code and set Data property
//
//     return nil, jsonapi.E301.SetData("http://google.com")
//
// Redirecting depends on http.Redirect(). The data returned from handler will never
// write to ResponseWriter.
//
// This basically obey the http://jsonapi.org rules:
//
//     - Return {"data": your_data} if error == nil
//     - Return {"errors": [{"code": application-defined-error-code, "detail": message}]} if error returned
type Handler func(r Request) (interface{}, error)

// ServeHTTP implements net/http.Handler
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	res, err := h(FromHTTP(w, r))
	resp := make(map[string]interface{})
	if err == nil {
		resp["data"] = res
		e := enc.Encode(resp)
		if e == nil {
			return
		}
		delete(resp, "data")

		err = E500.SetOrigin(e).SetData(
			`Failed to marshal data`,
		)
	}

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
	resp["errors"] = []*ErrObj{&ErrObj{Detail: err.Error()}}
	enc.Encode(resp)
}
