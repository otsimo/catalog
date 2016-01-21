package models

import "errors"

var (
	ErrorNotFound       = errors.New("Not found")
	ErrorNotImplemented = errors.New("Not implemented yet")
	ErrorInvalidID      = errors.New("Invalid ID")
	ErrorIDExists       = errors.New("ID already exists")
)
