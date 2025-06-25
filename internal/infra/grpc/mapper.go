package grpc

import (
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func mapProtoPackageType(pt pkg.PackageType) (*pvz.PackageType, error) {
	switch pt {
	case pkg.PackageType_PACKAGE_TYPE_UNSPECIFIED:
		return nil, nil
	case pkg.PackageType_PACKAGE_TYPE_BAG:
		t := pvz.PackageTypeBag
		return &t, nil
	case pkg.PackageType_PACKAGE_TYPE_BOX:
		t := pvz.PackageTypeBox
		return &t, nil
	case pkg.PackageType_PACKAGE_TYPE_TAPE:
		t := pvz.PackageTypeTape
		return &t, nil
	case pkg.PackageType_PACKAGE_TYPE_BAG_TAPE:
		t := pvz.PackageTypeBagTape
		return &t, nil
	case pkg.PackageType_PACKAGE_TYPE_BOX_TAPE:
		t := pvz.PackageTypeBoxTape
		return &t, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Ошибка при выборе упаковки Package Type")
	}
}

func mapAcceptOrderRequestToInput(req *pkg.AcceptOrderRequest) (pvz.OrderCreate, error) {
	var expiresAt time.Time
	if req.ExpiresAt != nil {
		expiresAt = req.ExpiresAt.AsTime()
	}

	pkgType, err := mapProtoPackageType(req.GetPackage())
	if err != nil {
		return pvz.OrderCreate{}, err
	}

	return pvz.OrderCreate{
		OrderID:   pvz.OrderID(req.OrderId),
		UserID:    pvz.UserID(req.UserId),
		ExpiresAt: expiresAt,
		Package:   pkgType,
		Weight:    req.Weight,
		Price:     decimal.NewFromFloat(float64(req.Price)),
	}, nil
}

func mapProcessOrdersRequestToInput(req *pkg.ProcessOrdersRequest) app.ProcessOrdersInput {
	return app.ProcessOrdersInput{
		UserID:   req.UserId,
		Action:   pvz.ActionType(req.Action),
		OrderIDs: req.OrderIds,
	}
}
func mapProcessResultToProto(res app.ProcessOrdersOutput) *pkg.ProcessResult {
	return &pkg.ProcessResult{
		Processed: res.Processed,
		Errors:    res.Errors,
	}
}

func mapDomainOrdersToProto(domainOrders []pvz.Order) []*pkg.Order {
	if domainOrders == nil {
		return nil
	}

	protoOrders := make([]*pkg.Order, 0, len(domainOrders))
	for _, o := range domainOrders {
		
		pt := pkg.PackageType_PACKAGE_TYPE_UNSPECIFIED
		if o.Package != nil {
			pt = mapDomainPackageToProto(*o.Package)
		}
		protoOrders = append(protoOrders, &pkg.Order{
			OrderId:    uint64(o.OrderID),
			UserId:     uint64(o.UserID),
			Status:     pkg.OrderStatus(o.Status),
			ExpiresAt:  timestamppb.New(o.ExpiresAt),
			Weight:     o.Weight,
			TotalPrice: float32(o.Price.InexactFloat64()),
			Package:    pt,
		})
	}

	return protoOrders
}

func mapDomainPackageToProto(pt pvz.PackageType) pkg.PackageType {
	switch pt {
	case pvz.PackageTypeBag:
		return pkg.PackageType_PACKAGE_TYPE_BAG
	case pvz.PackageTypeBox:
		return pkg.PackageType_PACKAGE_TYPE_BOX
	case pvz.PackageTypeTape:
		return pkg.PackageType_PACKAGE_TYPE_TAPE
	case pvz.PackageTypeBagTape:
		return pkg.PackageType_PACKAGE_TYPE_BAG_TAPE
	case pvz.PackageTypeBoxTape:
		return pkg.PackageType_PACKAGE_TYPE_BOX_TAPE
	default:
		return pkg.PackageType_PACKAGE_TYPE_UNSPECIFIED
	}
}

func mapDomainHistoryToProto(his app.OrderHistory) *pkg.OrderHistory {
	return &pkg.OrderHistory{
		OrderId:   his.OrderID,
		Status:    pkg.OrderStatus(his.Status),
		CreatedAt: timestamppb.New(his.Time),
	}
}

func mapDomainReturnsListToProto(returnList *app.ReturnsListOutput) *pkg.ReturnsList {
	if returnList == nil {
		return nil
	}
	return &pkg.ReturnsList{
		Returns: mapDomainOrdersToProto(returnList.Returns),
	}
}
