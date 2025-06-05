package pvz

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const DatetimeFormat = "2006-01-02"

type OrdersService struct {
	storage map[uint64]*pvz.Order
}

func NewOrdersService() *OrdersService {
	return &OrdersService{
		storage: make(map[uint64]*pvz.Order),
	}
}

func FormatTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	time := ts.AsTime()
	return time.Format(DatetimeFormat)
}

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
func packageTypeToString(pkg pvz.PackageType) string {
	switch pkg {
	case pvz.PackageType_PACKAGE_TYPE_UNSPECIFIED:
		return "UNSPECIFIED"
	case pvz.PackageType_PACKAGE_TYPE_BAG:
		return "BAG"
	case pvz.PackageType_PACKAGE_TYPE_BOX:
		return "BOX"
	case pvz.PackageType_PACKAGE_TYPE_TAPE:
		return "TAPE"
	case pvz.PackageType_PACKAGE_TYPE_BAG_TAPE:
		return "BAG_TAPE"
	case pvz.PackageType_PACKAGE_TYPE_BOX_TAPE:
		return "BOX_TAPE"
	default:
		return "UNKNOWN"
	}
}
