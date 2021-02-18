// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package apitool

import (
	"fmt"
	"testing"
)

func TestHOTPGen(t *testing.T) {
	secret := [10]byte{
		0xca, 0xfe, 0xba, 0xbe,
		0xde, 0xad, 0xbe, 0xef,
		0x4b, 0x1d,
	}

	// test data are generated from google authenticator debugging tool
	// see github.com/google/google-authenticator-libpam/blob/master/totp.html
	cases := []struct {
		value int64
		code  string
	}{
		{0, "323633"},
		{1, "178548"},
		{16, "000635"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("HOTP#%d", c.value), func(t *testing.T) {
			m := TOTPMiddleware{
				Secret: secret,
				Digit:  6,
			}

			if x := m.HOTP(c.value); x != c.code {
				t.Fatalf("expected %s, got %s", c.code, x)
			}
		})
	}
}
