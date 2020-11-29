package Week02

import (
	"database/sql"

	"github.com/pkg/errors"
)

func SetUser() error {
	return errors.Wrap(sql.ErrNoRows, "set user failed")
}

func IsErrNoRows(err error) bool {
	if errors.Cause(err) == sql.ErrNoRows {
		return true
	}
	return false
}
