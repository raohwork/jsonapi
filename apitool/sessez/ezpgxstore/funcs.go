// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ezpgxstore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/raohwork/jsonapi/apitool/sessez"
)

func actx() (ret context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func (s *store) db() (conn *pgxpool.Conn, err error) {
	ctx, cancel := actx()
	defer cancel()
	return s.pool.Acquire(ctx)
}

func (s *store) New(ttl time.Duration) (id string, err error) {
	conn, err := s.db()
	if err != nil {
		return
	}
	defer conn.Release()

	ctx, cancel := actx()
	defer cancel()
	row := conn.QueryRow(ctx, s.qNew, ttl/time.Millisecond)
	err = row.Scan(&id)
	return
}

func (s *store) Get(id string) (data map[string]sessez.InternalData, err error) {
	conn, err := s.db()
	if err != nil {
		return
	}
	defer conn.Release()

	ctx, cancel := actx()
	defer cancel()
	row := conn.QueryRow(ctx, s.qGet, id)
	err = row.Scan(&data)
	return
}

func (s *store) Set(id string, data map[string]sessez.InternalData) (err error) {
	conn, err := s.db()
	if err != nil {
		return
	}
	defer conn.Release()

	ctx, cancel := actx()
	defer cancel()
	res, err := conn.Exec(ctx, s.qSet, data, id)
	if err != nil {
		return
	}
	if res.RowsAffected() < 1 {
		return pgx.ErrNoRows
	}
	return
}

func (s *store) Unset(id string) (err error) {
	conn, err := s.db()
	if err != nil {
		return
	}
	defer conn.Release()

	ctx, cancel := actx()
	defer cancel()
	_, err = conn.Exec(ctx, s.qUnset, id)
	return
}

func (s *store) GC() {
	conn, err := s.db()
	if err != nil {
		return
	}
	defer conn.Release()

	ctx, cancel := actx()
	defer cancel()
	conn.Exec(ctx, s.qGC)
	return
}
