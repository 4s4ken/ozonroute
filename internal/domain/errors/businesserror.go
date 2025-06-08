package errors

import (
	"fmt"
)

type ErrorCode int64

const (
	ErrorCodeNotFound ErrorCode = iota + 1
	ErrorCodeAlreadyExists
	ErrorCodeExpires
	ErrorCodeValidation
	ErrorCodeInternal
	ErrorCodePackage
	ErrorCodeWeight
)

type Error struct {
	Message string
	Code    ErrorCode
}

func (e Error) Error() string {
	return e.Message
}

// EntityNotFoundError создает ошибку для случая, когда сущность не найдена.
func EntityNotFoundError(entityName string, entityID string) error {
	return &Error{
		Message: fmt.Sprintf("%s с идентификатором %s не найдена", entityName, entityID),
		Code:    ErrorCodeNotFound,
	}
}

// EntityAlreadyExistsError создает ошибку для случая, когда сущность уже существует.
func EntityAlreadyExistsError(entityName string, entityID string) error {
	return &Error{
		Message: fmt.Sprintf("%s с идентификатором %s уже существует", entityName, entityID),
		Code:    ErrorCodeAlreadyExists,
	}
}

// EntityExpiredError создает ошибку для случая, когда сущность просрочена.
func EntityExpiredError(entityName string, entityID string) error {
	return &Error{
		Message: fmt.Sprintf("%s с идентификатором %s просрочена", entityName, entityID),
		Code:    ErrorCodeExpires,
	}
}

// EntityValidationError создает ошибку для случая ошибки валидации сущности.
func EntityValidationError(entityName string, entityID string) error {
	return &Error{
		Message: fmt.Sprintf("Ошибка валидации %s с идентификатором %s", entityName, entityID),
		Code:    ErrorCodeValidation,
	}
}

// EntityInternalError создает ошибку для внутренней ошибки при обработке сущности.
func EntityInternalError(entityName string, entityID string) error {
	return &Error{
		Message: fmt.Sprintf("Внутренняя ошибка при обработке %s с идентификатором %s", entityName, entityID),
		Code:    ErrorCodeInternal,
	}
}

// EntityPackageError создает ошибку при выборе упаковки для заказа.
func EntityPackageError(entityName string, entityID string) error {
	return &Error{
		Message: fmt.Sprintf("Ошибка при выборе упаковки %s для заказа с идентификатором %s", entityName, entityID),
		Code:    ErrorCodePackage,
	}
}

// EntityWeightError создает ошибку, связанную с весом заказа и упаковкой.
func EntityWeightError(entityName string, entityID string) error {
	return &Error{
		Message: fmt.Sprintf("Слишком большой вес у заказа с идентификатором %s для упаковки %s", entityID, entityName),
		Code:    ErrorCodeWeight,
	}
}
