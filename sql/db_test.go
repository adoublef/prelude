package sql_test

import (
	"context"
	"embed"
	"fmt"
	"testing"

	"github.com/adoublef/prelude/sql"
	st "github.com/adoublef/prelude/sql/sqltest"
	"github.com/matryer/is"
)

//go:embed all:*.up.sql
var fsys embed.FS

func TestSQL(t *testing.T) {
	t.Run("query", st.RoundTrip(fsys, func(t *testing.T, db *sql.DB) {
		var (
			is  = is.New(t)
			ctx = context.Background()
			err error
		)

		err = db.QueryRow(ctx, "SELECT * FROM t WHERE t.a = 1").Err()
		is.NoErr(err) // exists

		var n int
		err = db.QueryRow(ctx, "SELECT * FROM t WHERE t.a = 2").Scan(&n)
		is.Equal(err, sql.ErrNoRows) // not exist

		var ok bool
		err = db.QueryRow(ctx, "SELECT EXISTS (SELECT * FROM t WHERE t.a = 1)").Scan(&ok)
		is.NoErr(err)      // query
		is.Equal(ok, true) // exists

		var bad bool
		err = db.QueryRow(ctx, fmt.Sprintf("SELECT EXISTS (%s)", "SELECT * FROM t WHERE t.a = 2")).Scan(&bad)
		is.NoErr(err)        // query
		is.Equal(bad, false) // does not exists
	}))
}
