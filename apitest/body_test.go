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
	"bytes"
	"testing"

	"github.com/raohwork/jsonapi"
)

func TestWithOrder(t *testing.T) {
	buf := &bytes.Buffer{}

	m1 := func(h jsonapi.Handler) jsonapi.Handler {
		return func(r jsonapi.Request) (interface{}, error) {
			buf.WriteByte('1')
			data, err := h(r)
			buf.WriteByte('1')
			return data, err
		}
	}
	m2 := func(h jsonapi.Handler) jsonapi.Handler {
		return func(r jsonapi.Request) (interface{}, error) {
			buf.WriteByte('2')
			data, err := h(r)
			buf.WriteByte('2')
			return data, err
		}
	}

	h := func(r jsonapi.Request) (interface{}, error) {
		buf.WriteByte('3')
		return nil, nil
	}

	Test(h).With(m1).With(m2).Use(nil)

	if actual := buf.String(); actual != "21312" {
		t.Fatalf("expected 21312, got %s", actual)
	}
}
