package storage_errors

import "errors"

var (
	INVALID_PACKAGE      = errors.New("") //???
	WEIGHT_TOO_HEAVY     = errors.New("") //???
	ORDER_NOT_FOUND      = errors.New("order not found")
	ORDER_ALREADY_EXISTS = errors.New("order already exists")
	STORAGE_EXPIRED      = errors.New("storage is expired")
	VALIDATION_FAILED    = errors.New("validation failed")
	INTERNAL_ERROR       = errors.New("") //???
)
