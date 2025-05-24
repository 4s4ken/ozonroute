package businesserror

import (
	"fmt"
)

type ErrorCode int64

const (
	ErrorCodeNotFound      ErrorCode = 1
	ErrorCodeAlreadyExists ErrorCode = 2
	ErrorCodeExpires       ErrorCode = 3
	ErrorCodeValidation    ErrorCode = 4
	ErrorCodeInternal      ErrorCode = 5
)

type Error struct {
	Message string
	Code    ErrorCode
}

func (e Error) Error() string {
	return e.Message
}

func EntityNotFoundError(entityName string, entityID string) error {
	return Error{
		Message: fmt.Sprintf("%s с идентификатором %s не найдена", entityName, entityID),
		Code:    ErrorCodeNotFound,
	}
}

func EntityAlreadyExistsError(entityName string, entityID string) error {
	return Error{
		Message: fmt.Sprintf("%s с идентификатором %s уже существует", entityName, entityID),
		Code:    ErrorCodeAlreadyExists,
	}
}
func EntityExpiredError(entityName string, entityID string) error {
	return Error{
		Message: fmt.Sprintf("%s с идентификатором %s просрочена", entityName, entityID),
		Code:    ErrorCodeExpires,
	}
}
func EntityValidationError(entityName string, entityID string) error {
	return Error{
		Message: fmt.Sprintf("Ошибка валидации %s с идентификатором %s", entityName, entityID),
		Code:    ErrorCodeValidation,
	}
}
func EntityInternalError(entityName string, entityID string) error {
	return Error{
		Message: fmt.Sprintf("Внутренняя ошибка при обработке %s с идентификатором %s", entityName, entityID),
		Code:    ErrorCodeInternal,
	}
}
