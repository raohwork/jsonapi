// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

// EFormat indicates an error when parsing server response.
//
// In general, EFormat represents errors occurred after request is sent.
type EFormat struct {
	Origin error
}

func (e EFormat) Error() string {
	ret := "cannot parse response from api server"
	if e.Origin != nil {
		ret += ": " + e.Origin.Error()
	}
	return ret
}

func (e EFormat) Unwrap() error { return e.Origin }

// EClient indicates this api call is failed due to client side error.
//
// In general, Eclient represents errors occurred before sending request.
type EClient struct {
	Origin error
}

func (e EClient) Error() string {
	ret := "there's something wrong at client side"
	if e.Origin != nil {
		ret += ": " + e.Origin.Error()
	}
	return ret
}

func (e EClient) Unwrap() error { return e.Origin }
