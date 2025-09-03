package utils

import "errors"

var (
	ErrValidation = errors.New("validation failed")
	ErrConflict   = errors.New("conflict")
	ErrInternal   = errors.New("internal error")
)
