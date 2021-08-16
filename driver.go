package sqljs

import (
	"database/sql"
	"database/sql/driver"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"
)

func init() {
	sql.Register("sqljs", &Driver{})
}

type Driver struct {
	LocateFile func(string) string
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	SQL, err := promise.Await(js.Global().Call("initSqlJs", d.config()))
	if err != nil {
		return nil, err
	}
	// FIXME name could be base64 encoded data...
	db := SQL.Get("Database").New()
	return &Conn{db}, nil
}

func (d *Driver) config() map[string]interface{} {
	config := make(map[string]interface{}, 1)
	if d.LocateFile != nil {
		config["locateFile"] = js.FuncOf(d.LocateFile)
	}
	return config
}
