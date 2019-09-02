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
