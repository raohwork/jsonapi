// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package sessez

import (
	"time"
)

// Store defines an interface to persist/load data in storage.
type Store interface {
	// Creates a new session id, failed if internal error.
	New(ttl time.Duration) (id string, err error)
	// Get session data, nil if not found, failed if internal error.
	//
	// Internal TTL info *SHOULD* be updated.
	Get(id string) (data map[string]InternalData, err error)
	// Update session data for id, failed if ErrSessionNotFound or internal
	// error.
	//
	// Internal TTL info *MUST* be updated.
	Set(id string, data map[string]InternalData) error
	// Deletes a session id and related data, failed if ErrSessionNotFound or
	// internal error.
	Unset(id string) error
	// Requests to delete outdated session data.
	// It should be a NO-OP if the store does not support it.
	GC()
}
