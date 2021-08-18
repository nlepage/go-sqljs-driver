package sqljs

import (
	"context"
	"database/sql/driver"
	"errors"
	"syscall/js"
)

type Conn struct {
	js.Value
}

func (c *Conn) Prepare(query string) (stmt driver.Stmt, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoveredToError(r)
		}
	}()

	stmt = &Stmt{c.Call("prepare", query)}

	return
}

func (c *Conn) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoveredToError(r)
		}
	}()

	c.Call("close")

	return
}

func (c *Conn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

var _ driver.ConnBeginTx = &Conn{}

func (c *Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

// FIXME All Conn implementations should implement the following interfaces: Pinger, SessionResetter, and Validator.

// FIXME If named parameters or context are supported, the driver's Conn should implement: ExecerContext, QueryerContext, ConnPrepareContext, and ConnBeginTx.
