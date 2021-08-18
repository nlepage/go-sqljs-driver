package sqljs

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
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
			err = errors.New(r.(string))
		}
	}()

	if ok := stmt.Call("run", namedValuesToBindParams(args)).Bool(); !ok {
		err = errors.New("statement was not reset") // return value of run is the return value of reset
		return
	}

	res = &Result{}

	return
}

var _ driver.StmtQueryContext = &Stmt{}

func (stmt *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(js.Error)
		}
	}()

	// Don't check return value, bind always returns true
	stmt.Call("bind", namedValuesToBindParams(args))

	rows = &StmtRows{stmt.Value}

	return
}

type StmtRows struct {
	js.Value
}

func (stmt *StmtRows) Columns() []string {
	v := stmt.Call("getColumnNames")
	length := v.Length()
	columns := make([]string, length)
	for i := 0; i < length; i++ {
		columns[i] = v.Index(i).String()
	}
	return columns
}

func (stmt *StmtRows) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(js.Error)
		}
	}()

	if ok := stmt.Call("reset").Bool(); !ok {
		err = errors.New("statement was not reset")
		return
	}

	return
}

func (stmt *StmtRows) Next(dest []driver.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(js.Error)
		}
	}()

	if ok := stmt.Call("step").Bool(); !ok {
		err = io.EOF
		return
	}

	row := stmt.Call("get")
	length := row.Length()
	for i := 0; i < length; i++ {
		v := row.Index(i)
		switch v.Type() {
		case js.TypeNull:
			dest[i] = nil
		case js.TypeNumber:
			dest[i] = v.Float()
		case js.TypeString:
			dest[i] = v.String()
		// FIXME Uint8Array to []byte
		default:
			err = fmt.Errorf("unknown SqlValue type %s", v.Type())
			return
		}
	}

	return
}

func valuesToNamedValues(values []driver.Value) []driver.NamedValue {
	nValues := make([]driver.NamedValue, len(values))
	for i, value := range values {
		nValues[i].Ordinal = i + 1
		nValues[i].Value = value
	}
	return nValues
}

// see https://sql.js.org/documentation/Statement.html#.BindParams
func namedValuesToBindParams(values []driver.NamedValue) interface{} {
	if values == nil {
		return nil
	}

	named := false
	for _, value := range values {
		if value.Name != "" {
			named = true
			break
		}
	}

	if named {
		params := make(map[string]interface{}, len(values))
		for _, value := range values {
			if value.Name == "" {
				params[strconv.Itoa(value.Ordinal)] = value.Value
			} else {
				params[value.Name] = value.Value
			}
		}
		return params
	}

	params := make([]interface{}, len(values))
	for _, value := range values {
		params[value.Ordinal-1] = value.Value
	}
	return params
}
