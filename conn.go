package sqljs

import (
	"database/sql/driver"
	"errors"
	"syscall/js"
)

type Conn struct {
	js.Value
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (c *Conn) Close() error {
	return errors.New("NOT IMPLEMENTED")
}

// Deprecated: Drivers should implement ConnBeginTx instead (or additionally).
func (c *Conn) Begin() (driver.Tx, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}
