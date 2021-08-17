package sqljs

import "errors"

type Result struct{}

func (r *Result) LastInsertId() (int64, error) {
	return 0, errors.New("not implemented")
}

func (r *Result) RowsAffected() (int64, error) {
	return 0, errors.New("not implemented")
}
