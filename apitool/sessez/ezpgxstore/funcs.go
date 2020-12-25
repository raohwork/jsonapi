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
