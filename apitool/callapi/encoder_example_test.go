// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

func MyEncoder(v any) ([]byte, error) {
	switch x := v.(type) {
	case map[string][]string:
		return []byte(url.Values(x).Encode()), nil
	case map[string]string:
		val := url.Values{}
		for k, v := range x {
			val.Set(k, v)
		}
		return []byte(val.Encode()), nil
	}
	return nil, errors.New("unsupported type")
}

type myAPIResp struct {
	Status string          `json:"status"` // ok or fail
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
}

func MyParser(resp *http.Response, result interface{}) (err error) {
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	var tmp myAPIResp
	err = json.NewDecoder(resp.Body).Decode(&tmp)
	if err != nil {
		return
	}

	if tmp.Status == "ok" {
		return json.Unmarshal(tmp.Data, result)
	}

	if tmp.Error != "" {
		return errors.New(tmp.Error)
	}

	return errors.New("unknown error returned from server")
}

func Example_customEncoder() {
	// this endpoint accepts post form, and returns json
	ep := Encoder(MyEncoder).
		EP(http.MethodPost, "https://example.com/api/my_endpoint").
		SendBy(nil).
		ParseWith(MyParser)
	// var result MyResult
	// err = ep.Call(param, &result)

	_ = ep
}
