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

package apitest

import (
	"testing"

	"github.com/raohwork/jsonapi"
)

// AssertError validates if returned values is specified error value
//
// It checks for following situations:
//
//     - data must be nil
//     - err must be jsonapi.Error type
//     - expect.EqualTo(err) == true
func AssertError(t *testing.T, expect jsonapi.Error, data interface{}, err error) jsonapi.Error {
	if data != nil {
		t.Errorf("handler in error state should not return any data, got %#v", data)
	}

	if err == nil {
		t.Fatal("handler in error state should return a error")
	}

	e, ok := err.(jsonapi.Error)
	if !ok {
		t.Fatalf("handler in error state should return api errors, got %#v", err)
	}

	if !expect.EqualTo(e) {
		t.Errorf("error should be %s, got %s", expect, e)
	}

	return e
}
