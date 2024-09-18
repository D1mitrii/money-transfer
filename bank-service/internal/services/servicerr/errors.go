package servicerr

import "errors"

var (
	ErrAlreadyExist    = errors.New("already exist")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
)
