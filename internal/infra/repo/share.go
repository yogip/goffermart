package repo

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"syscall"

	"github.com/jackc/pgx"
)

var recoverableErrors = []error{
	syscall.ECONNREFUSED,
	pgx.ErrDeadConn,
	sql.ErrConnDone,
	pgx.ErrConnBusy,
	driver.ErrBadConn,
	io.EOF,
}
