// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

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
