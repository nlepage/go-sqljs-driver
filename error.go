package sqljs

import (
	"errors"
	"fmt"
)

func recoveredToError(r interface{}) error {
	switch err := r.(type) {
	case error:
		return err
	case string:
		return errors.New(err)
	case fmt.Stringer:
		return errors.New(err.String())
	default:
		return fmt.Errorf("%v", err)
	}
}
