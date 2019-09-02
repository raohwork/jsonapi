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

package jsonapi

import (
	"context"
	"encoding/json"
	"net/http"
)

// Request represents most used data a handler need
type Request interface {
	// Decode() helps you to read parameters in request body
	Decode(interface{}) error
	// R() retrieves original http request
	R() *http.Request
	// W() retrieves original http response writer
	W() http.ResponseWriter
	// WithValue() adds a new key-value pair in context of http request
	WithValue(key, val interface{}) Request
}

// FakeRequest implements a Request and let you do some magic in it
type FakeRequest struct {
	// this is used to implement Request.Decode()
	Decoder *json.Decoder
	// this is used to implement Request.R() and Request.WithValue()
	Req *http.Request
	// this is used to implement Request.W()
	Resp http.ResponseWriter
}

// Decode implements Request
func (r *FakeRequest) Decode(data interface{}) error {
	return r.Decoder.Decode(data)
}

// R implements Request
func (r *FakeRequest) R() *http.Request {
	return r.Req
}

// W implements Request
func (r *FakeRequest) W() http.ResponseWriter {
	return r.Resp
}

// WithValue implements Request
func (r *FakeRequest) WithValue(key, val interface{}) (ret Request) {
	req := r.Req
	req = req.WithContext(
		context.WithValue(
			req.Context(),
			key,
			val,
		),
	)
	return &FakeRequest{
		Decoder: r.Decoder,
		Req:     req,
		Resp:    r.Resp,
	}
}

// FromHTTP creates a Request instance from http request and response
func FromHTTP(w http.ResponseWriter, r *http.Request) Request {
	dec := json.NewDecoder(r.Body)
	return &FakeRequest{
		Decoder: dec,
		Req:     r,
		Resp:    w,
	}
}

type reqWrapper struct {
	Request
	req *http.Request
}

func (r *reqWrapper) R() *http.Request {
	return r.req
}

// WrapRequest creates a new Request, with http request replaced
func WrapRequest(q Request, r *http.Request) Request {
	return &reqWrapper{
		Request: q,
		req:     r,
	}
}

type respWrapper struct {
	Request
	resp http.ResponseWriter
}

func (r *respWrapper) W() http.ResponseWriter {
	return r.resp
}

// WrapResponse creates a new Request, with http response replaced
func WrapResponse(q Request, w http.ResponseWriter) Request {
	return &respWrapper{
		Request: q,
		resp:    w,
	}
}
