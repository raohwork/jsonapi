// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

import (
	"encoding/json"
	"net/http"

	"github.com/raohwork/jsonapi"
)

// Parser parses server response and decode it. It takes responsibility to close
// response body.
type Parser func(resp *http.Response, result interface{}) error

type callResp struct {
	Data   *json.RawMessage `json:"data"`
	Errors []jsonapi.ErrObj `json:"errors"`
}

// DefaultParser parses response of a jsonapi
//
// If any io or json parsing error occurred, an EFormat is returned.
func DefaultParser(resp *http.Response, result interface{}) error {
	var res callResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return EFormat{err}
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
