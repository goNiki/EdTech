package errorsAPP

import "errors"

var (
	ErrFailValidate       = errors.New("failed validate fields")
	ErrFailHashingPassword = errors.New("fieled hash password")
)
