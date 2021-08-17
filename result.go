package sqljs

import "errors"

type Result struct{}

func (res *Result) LastInsertId() (int64, error) {
	return 0, errors.New("not implemented")
}

func (res *Result) RowsAffected() (int64, error) {
	return 0, errors.New("not implemented")
}
