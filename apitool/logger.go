// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package apitool

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/raohwork/jsonapi"
)

// SimpleFormat creates a log provider logs only error messages
//
// It logs original error instead if exists.
func SimpleFormat(l *log.Logger) LogProvider {
	return LogProvider(func(r *http.Request, data interface{}, err error) {
		if e, ok := err.(jsonapi.Error); ok && e.Origin != nil {
			err = e.Origin
		}
		l.Print(err)
	})
}

// BasicFormat creates a log provider logs url and error messages
//
// It logs original error instead if exists.
func BasicFormat(l *log.Logger) LogProvider {
	return LogProvider(func(r *http.Request, data interface{}, err error) {
		if e, ok := err.(jsonapi.Error); ok && e.Origin != nil {
			err = e.Origin
		}
		l.Printf("%s: %s", r.URL, err)
	})
}

// JSONLog represents all entries a JSON log provider will log
type JSONLog struct {
	Method     string         `json:"request_method"`
	URL        *url.URL       `json:"request_url"`
	Header     http.Header    `json:"request_header"`
	Host       string         `json:"request_host"`
	RemoteAddr string         `json:"request_remote_addr"`
	Cookies    []*http.Cookie `json:"cookies"`
	Data       interface{}    `json:"reply_data"`
	Error      error          `json:"reply_error"`
}

// JSONFormat creates a log provider logs detailed info in json format
func JSONFormat(l *log.Logger) LogProvider {
	return LogProvider(func(r *http.Request, data interface{}, err error) {
		buf, _ := json.Marshal(&JSONLog{
			Method:     r.Method,
			URL:        r.URL,
			Header:     r.Header,
			Host:       r.Host,
			RemoteAddr: r.RemoteAddr,
			Cookies:    r.Cookies(),
			Data:       data,
			Error:      err,
		})
		l.Print(string(buf))
	})
}
