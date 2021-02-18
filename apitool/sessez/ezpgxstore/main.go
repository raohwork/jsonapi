// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ezpgxstore

import (
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/raohwork/jsonapi/apitool/sessez"
)

type store struct {
	pool *pgxpool.Pool

	qNew   string
	qGet   string
	qSet   string
	qUnset string
	qGC    string
}

// CreateTable tries to create table with minimum data structure.
func CreateTable(table string, pool *pgxpool.Pool) (err error) {
	ctx, cancel := actx()
	defer cancel()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	qext := `CREATE EXTENSION IF NOT EXISTS pgcrypto`
	ctx, cancel = actx()
	defer cancel()
	if _, err = conn.Exec(ctx, qext); err != nil {
		return
	}

	qtbl := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
  id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
  data json NOT NULL,
  ttl integer NOT NULL,
  last_used timestamp with time zone NOT NULL
)`, table)
	ctx, cancel = actx()
	defer cancel()
	if _, err = conn.Exec(ctx, qtbl); err != nil {
		return
	}

	qidx := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s_last_used_idx ON %s USING btree (last_used ASC NULLS LAST)`, table, table)
	ctx, cancel = actx()
	defer cancel()
	if _, err = conn.Exec(ctx, qidx); err != nil {
		return
	}

	return
}

// New creates a session store which persists session data in PostgreSQL using pgx
func New(table string, pool *pgxpool.Pool) (ret sessez.Store) {
	return &store{
		pool: pool,
		qNew: fmt.Sprintf(
			`INSERT INTO %s ("data", "ttl", "last_used") VALUES ('{}', $1, CURRENT_TIMESTAMP) RETURNING "id"`,
			table,
		),
		qGet: fmt.Sprintf(
			`UPDATE %s SET "last_used"=CURRENT_TIMESTAMP
WHERE "id"=$1
  AND "last_used" >= CURRENT_TIMESTAMP - ("ttl"::TEXT||' milliseconds')::interval
RETURNING "data"`,
			table,
		),
		qSet: fmt.Sprintf(
			`UPDATE %s SET "data"=$1, "last_used"=CURRENT_TIMESTAMP
WHERE "id"=$2
  AND "last_used" >= CURRENT_TIMESTAMP - ("ttl"::TEXT||' milliseconds')::interval`,
			table,
		),
		qUnset: fmt.Sprintf(
			`DELETE FROM %s WHERE "id"=$1`, table,
		),
		qGC: fmt.Sprintf(
			`DELETE FROM %s WHERE age("last_used") > ("ttl"::TEXT||' seconds')::interval`,
			table,
		),
	}
}
