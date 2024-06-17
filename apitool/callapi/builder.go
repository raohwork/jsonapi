// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callapi

import (
	"net/http"
)

// Builder is a builder to build many [Caller] with same settings. Zero value builds
// same [Caller] as [EP] does.
//
// It is suggested to use Builder if any of following rules is matched:
//
//   - There're too many endpoints; using Builder can save some key strokes.
//   - Most of endpoints returns some data in response header.
//   - Most of endpoints need dynamic value in header, message digest and etag for example.
type Builder struct {
	// Maker is a function to create request. DefaultEncoder().EP is used if nil
	Maker func(method, url string) Endpoint
	// [http.DefaultClient] is used if nil
	Sender *http.Client
	// DefaultParser is used if nil
	Parser Parser
}

// EP creates a [Caller] that uses b.Maker to create request, send the request by
// b.Sender and parse the response by b.Parser.
func (b Builder) EP(method, uri string) Caller {
	if b.Maker == nil {
		b.Maker = DefaultEncoder().EP
	}
	if b.Parser == nil {
		b.Parser = DefaultParser
	}
	return b.Maker(method, uri).SendBy(b.Sender).ParseWith(b.Parser)
}

// UseMaker creates a new Builder that use m as maker.
func (b Builder) UseMaker(m func(method, url string) Endpoint) Builder {
	b.Maker = m
	return b
}

// UseSender creates a new Builder that use cl as sender.
func (b Builder) UseSender(cl *http.Client) Builder {
	b.Sender = cl
	return b
}

// UseParser creates a new Builder that use p as parser.
func (b Builder) UseParser(p Parser) Builder {
	b.Parser = p
	return b
}
