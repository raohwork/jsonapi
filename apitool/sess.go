// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package apitool

import (
	"context"

	"github.com/Ronmi/rtoolkit/session"
	"github.com/raohwork/jsonapi"
)

// Session creates a api middleware that handles session related functions
//
// If you are facing "Trailer Header" problem with original session middleware,
// this should be helpful.
//
//	jsonapi.With(
//	    apitool.Session(mySessMgr),
//	).RegisterAll(myHandlerClass)
//
// Created middleware will try to save update cookie ttl value if possible. It
// fails silently.
//
// You have add build tag `rtoolkit_session` to use this middleware.
func Session(m *session.Manager) jsonapi.Middleware {
	return func(h jsonapi.Handler) jsonapi.Handler {
		return func(req jsonapi.Request) (i interface{}, e error) {
			r := req.R()
			sess, err := m.Start(req.W(), r)
			if err != nil {
				return nil, jsonapi.E500.SetOrigin(err)
			}

			r = r.WithContext(context.WithValue(
				r.Context(),
				session.SessionObjectKey,
				sess,
			))
			i, e = h(jsonapi.WrapRequest(req, r))

			_ = sess.Save(req.W())
			return
		}
	}
}
