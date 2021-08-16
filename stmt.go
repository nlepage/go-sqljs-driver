package sqljs

import (
	"context"
	"database/sql/driver"
	"errors"
	"syscall/js"
)

type Stmt struct {
	js.Value
}

// Deprecated: Drivers should implement StmtExecContext instead (or additionally).
func (stmt *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("not implemented")
}

// Deprecated: Drivers should implement StmtQueryContext instead (or additionally).
func (stmt *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, errors.New("not implemented")
}

func (stmt *Stmt) NumInput() int {
	return -1 // No equivalent on https://sql.js.org/documentation/Statement.html
}

func (stmt *Stmt) Close() error {
	if ok := stmt.Call("free").Bool(); !ok {
		return errors.New("statement did not close")
	}

	return nil
}

var _ driver.StmtExecContext = &Stmt{}

func (stmt *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return nil, errors.New("not implemented")
}

var _ driver.StmtQueryContext = &Stmt{}

func (stmt *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return nil, errors.New("not implemented")
}
