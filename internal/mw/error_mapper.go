package mw

import (
	"errors"
	businesserror "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MapBusinessErrorToGRPC конвертирует бизнес-ошибку в gRPC статус.
func MapBusinessErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	var be *businesserror.Error
	if errors.As(err, &be) {
		switch be.Code {
		case businesserror.ErrorCodeNotFound:
			return status.Errorf(codes.NotFound, be.Message)
		case businesserror.ErrorCodeAlreadyExists:
			return status.Errorf(codes.AlreadyExists, be.Message)
		case businesserror.ErrorCodeExpires:
			return status.Errorf(codes.FailedPrecondition, be.Message)
		case businesserror.ErrorCodeValidation:
			return status.Errorf(codes.InvalidArgument, be.Message)
		case businesserror.ErrorCodePackage, businesserror.ErrorCodeWeight:
			return status.Errorf(codes.InvalidArgument, be.Message)
		case businesserror.ErrorCodeInternal:
			return status.Errorf(codes.Internal, be.Message)
		default:
			return status.Errorf(codes.Unknown, be.Message)
		}
	}

	return status.Errorf(codes.Internal, err.Error())
}
