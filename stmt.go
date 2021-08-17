package sqljs

import (
	"context"
	"database/sql/driver"
	"errors"
	"strconv"
	"syscall/js"
)

type Stmt struct {
	js.Value
}

func (stmt *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	return stmt.ExecContext(context.Background(), valuesToNamedValues(args))
}

func (stmt *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return stmt.QueryContext(context.Background(), valuesToNamedValues(args))
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

func (stmt *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(js.Error)
		}
	}()

	stmt.Call("run", namedValuesToBindParams(args))

	res = &Result{}

	return
}

var _ driver.StmtQueryContext = &Stmt{}

func (stmt *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return nil, errors.New("not implemented")
}

func valuesToNamedValues(values []driver.Value) []driver.NamedValue {
	nValues := make([]driver.NamedValue, len(values))
	for i, value := range values {
		nValues[i].Ordinal = i + 1
		nValues[i].Value = value
	}
	return nValues
}

func namedValuesToBindParams(values []driver.NamedValue) interface{} {
	named := false
	for _, value := range values {
		if value.Name != "" {
			named = true
			break
		}
	}

	if named {
		params := make(map[string]driver.Value, len(values))
		for _, value := range values {
			if value.Name == "" {
				params[strconv.Itoa(value.Ordinal)] = value.Value
			} else {
				params[value.Name] = value.Value
			}
		}
		return params
	}

	params := make([]driver.Value, len(values))
	for _, value := range values {
		params[value.Ordinal-1] = value.Value
	}
	return params
}
