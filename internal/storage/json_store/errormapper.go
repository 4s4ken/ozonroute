package json_store

import (
	"errors"
	"fmt"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
)

func MapError(err error) error {
	var be *businesserror.Error
	if errors.As(err, &be) {
		switch be.Code {
		case businesserror.ErrorCodeNotFound:
			return notFoundError(be)
		case businesserror.ErrorCodeAlreadyExists:
			return alreadyExistsError(be)
		case businesserror.ErrorCodeExpires:
			return expiredError(be)
		case businesserror.ErrorCodeValidation:
			return validationError(be)
		case businesserror.ErrorCodeInternal:
			return internalError(be)
		default:
			return unknownError(be)
		}
	}
	return internalError(err)
}

func notFoundError(err error) error {
	return fmt.Errorf("NOT_FOUND: %w", err)
}

func alreadyExistsError(err error) error {
	return fmt.Errorf("ALREADY_EXISTS: %w", err)
}

func expiredError(err error) error {
	return fmt.Errorf("EXPIRED: %w", err)
}

func validationError(err error) error {
	return fmt.Errorf("VALIDATION_FAILED: %w", err)
}

func internalError(err error) error {
	return fmt.Errorf("INTERNAL_ERROR: %w", err)
}

func unknownError(err error) error {
	return fmt.Errorf("UNKNOWN_ERROR_CODE: %w", err)
}
