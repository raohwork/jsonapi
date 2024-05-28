// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gorsess

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/raohwork/jsonapi"
	"github.com/raohwork/jsonapi/apitool"
)

func testHandler(ctx context.Context, r jsonapi.Request) (data interface{}, err error) {
	sess, ok := apitool.GetSession(ctx, "reqkey")
	if !ok {
		err = errors.New("no session injected")
		return
	}

	r.R().ParseForm()
	val := r.R().Form.Get("param")
	if val != "" {
		sess.Set("data", val)
	}

	if r.R().Form.Get("discard") != "" {
		err = sess.Discard()
		data = err == nil
		return
	}

	if err = sess.Save(); err != nil {
		return
	}

	data, ok = sess.Get("data")
	data = map[string]interface{}{
		"data": data,
		"ok":   ok,
	}

	return
}

func TestSession(t *testing.T) {
	store := sessions.NewCookieStore(
		[]byte("aisuydfghjposmdg897ay8hid"),
		[]byte("123456789012345678901234"),
	)
	sp := New(store, "wtf")
	middleware := apitool.Session(sp, "reqkey")
	handler := middleware(testHandler)
	server := httptest.NewServer(handler)
	defer server.Close()
	s := testSuit{server, nil}

	ok := t.Run("setData", s.testSetData)
	if ok {
		ok = t.Run("getData", s.testGetData)
	}
	if ok {
		ok = t.Run("discard", s.testDiscard)
	}
}

type testSuit struct {
	*httptest.Server
	cookie *http.Cookie
}

func (s *testSuit) testSetData(t *testing.T) {
	hc := http.DefaultClient
	resp, err := hc.Get(s.URL + "?param=123")
	if err != nil {
		t.Fatal("sending test request: unexpected error: ", err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("reading response: unexpected error: ", err)
	}

	const expect = `{"data":{"data":"123","ok":true}}`
	str := strings.TrimSpace(string(buf))
	if str != expect {
		t.Log("expect output: ", expect)
		t.Log("actual output: ", str)
		t.Error("unexpected output")
	}

	ok := false
	for _, c := range resp.Cookies() {
		if c.Name == "wtf" {
			ok = true
			s.cookie = c
		}
	}
	if !ok {
		for _, c := range resp.Cookies() {
			t.Logf("all cookies: %+v", c)
		}
		t.Error("no cookie named 'wtf' found")
	}
}

func (s *testSuit) testGetData(t *testing.T) {
	hc := http.DefaultClient
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		t.Fatal("building test request: unexpected error: ", err)
	}
	req.AddCookie(s.cookie)

	resp, err := hc.Do(req)
	if err != nil {
		t.Fatal("sending test request: unexpected error: ", err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("reading response: unexpected error: ", err)
	}

	const expect = `{"data":{"data":"123","ok":true}}`
	str := strings.TrimSpace(string(buf))
	if str != expect {
		t.Log("expect output: ", expect)
		t.Log("actual output: ", str)
		t.Error("unexpected output")
	}

	ok := false
	for _, c := range resp.Cookies() {
		if c.Name == "wtf" {
			ok = true
			s.cookie = c
		}
	}
	if !ok {
		for _, c := range resp.Cookies() {
			t.Logf("all cookies: %+v", c)
		}
		t.Error("no cookie named 'wtf' found")
	}
}

func (s *testSuit) testDiscard(t *testing.T) {
	hc := http.DefaultClient
	req, err := http.NewRequest("GET", s.URL+"?discard=1", nil)
	if err != nil {
		t.Fatal("building test request: unexpected error: ", err)
	}
	req.AddCookie(s.cookie)

	resp, err := hc.Do(req)
	if err != nil {
		t.Fatal("sending test request: unexpected error: ", err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("reading response: unexpected error: ", err)
	}

	const expect = `{"data":true}`
	str := strings.TrimSpace(string(buf))
	if str != expect {
		t.Log("expect output: ", expect)
		t.Log("actual output: ", str)
		t.Error("unexpected output")
	}

	ok := true
	for _, c := range resp.Cookies() {
		if c.Name != "wtf" {
			continue
		}

		if c.Expires.After(time.Now()) {
			ok = false
			t.Logf("session cookie still there: %+v", c)
		}
		break
	}

	if !ok {
		t.Fatal("unexpected cookie")
	}
}
