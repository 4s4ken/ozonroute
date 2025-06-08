package grpc

import (
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const DatetimeFormat = "2006-01-02"

func FormatTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	time := ts.AsTime()
	return time.Format(DatetimeFormat)
}

func packageTypeToString(pkg pvz.PackageType) string {
	switch pkg {
	case pvz.PackageTypeUnspecified:
		return "UNSPECIFIED"
	case pvz.PackageTypeBag:
		return "BAG"
	case pvz.PackageTypeBox:
		return "BOX"
	case pvz.PackageTypeTape:
		return "TAPE"
	case pvz.PackageTypeBagTape:
		return "BAG_TAPE"
	case pvz.PackageTypeBoxTape:
		return "BOX_TAPE"
	default:
		return "UNKNOWN"
	}
}

func mapPackageProtoToDomain(p *pkg.PackageType) pvz.PackageType {
	if p == nil {
		return pvz.PackageTypeUnspecified
	}
	return pvz.PackageType(*p)
}

func mapAcceptOrderRequestToInput(req *pkg.AcceptOrderRequest) pvz.AcceptOrderInput {
	var expiresAt time.Time
	if req.ExpiresAt != nil {
		expiresAt = req.ExpiresAt.AsTime()
	}

	var pkgType *pvz.PackageType
	if req.Package != nil {
		pt := pvz.PackageType(*req.Package)
		pkgType = &pt
	}

	return pvz.AcceptOrderInput{
		OrderID:   req.OrderId,
		UserID:    req.UserId,
		ExpiresAt: expiresAt,
		Package:   pkgType,
		Weight:    req.Weight,
		Price:     req.Price,
	}
}

func mapReturnOrderRequestToInput(req *pkg.ReturnOrderRequest) pvz.ReturnOrderInput {
	return pvz.ReturnOrderInput{
		OrderID: pvz.OrderID(req.OrderId),
	}
}

func mapProcessOrdersRequestToInput(req *pkg.ProcessOrdersRequest) pvz.ProcessOrdersInput {
	return pvz.ProcessOrdersInput{
		UserID:   pvz.UserID(req.UserId),
		Action:   pvz.ActionType(req.Action),
		OrderIDs: pvz.OrderIDS{},
	}
}
func mapProcessResultToProto(res pvz.ProcessResult) *pkg.ProcessResult {
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
	for _, domainOrder := range domainOrders {
		protoOrders = append(protoOrders, &pkg.Order{
			OrderId: uint64(domainOrder.OrderID),
			Status:  pkg.OrderStatus(domainOrder.Status),
		})
	}

	return protoOrders
}

func mapDomainHistoryToProto(his pvz.OrderHistory) *pkg.OrderHistory {
	return &pkg.OrderHistory{
		OrderId:   his.OrderID,
		Status:    pkg.OrderStatus(his.Status),
		CreatedAt: timestamppb.New(his.CreatedAt),
	}
}

func mapDomainReturnsListToProto(returnList *pvz.ReturnsList) *pkg.ReturnsList {
	if returnList == nil {
		return nil
	}
	return &pkg.ReturnsList{
		Returns: mapDomainOrdersToProto(returnList.Returns),
	}
}
