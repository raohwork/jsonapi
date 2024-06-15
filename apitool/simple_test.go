// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package apitool

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/raohwork/jsonapi"
)

func TestSimpleLogger(t *testing.T) {
	cases := []struct {
		name   string
		e      error
		expect string
	}{
		{
			name:   "e500",
			e:      jsonapi.E500.SetData("test"),
			expect: "500: test",
		},
		{
			name:   "e500",
			e:      jsonapi.E500.SetOrigin(errors.New("test")),
			expect: "test",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			p := SimpleFormat(log.New(buf, "", 0))

			p(nil, nil, c.e)

			if actual := buf.String(); actual != c.expect+"\n" {
				t.Fatalf("expected [%s], got [%s]", c.expect, actual)
			}
		})
	}
}
