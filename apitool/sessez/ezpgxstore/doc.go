// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package ezpgxstore stores session data in Postgres using pgx
//
// You have to call GC() periodically to release db storage space.
package ezpgxstore
