package grpc

import (
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
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

func mapPackageProtoToDomain(pac *pkg.PackageType) pvz.PackageType {
	if pac == nil {
		return "nil"
	}
	return pvz.PackageType(*pac)
}

func mapAcceptOrderRequestToInput(req *pkg.AcceptOrderRequest) pvz.OrderCreate {
	var expiresAt time.Time
	if req.ExpiresAt != nil {
		expiresAt = req.ExpiresAt.AsTime()
	}

	var pkgType *pvz.PackageType
	if req.Package != nil {
		pt := pvz.PackageType(*req.Package)
		pkgType = &pt
	}

	return pvz.OrderCreate{
		OrderID:   pvz.OrderID(req.OrderId),
		UserID:    pvz.UserID(req.UserId),
		ExpiresAt: expiresAt,
		Package:   pkgType,
		Weight:    req.Weight,
		Price:     req.Price,
	}
}

func mapReturnOrderRequestToInput(req *pkg.ReturnOrderRequest) pvz.OrderCreate {
	return pvz.OrderCreate{
		OrderID: pvz.OrderID(req.OrderId),
	}
}

func mapProcessOrdersRequestToInput(req *pkg.ProcessOrdersRequest) app.ProcessOrdersInput {
	return app.ProcessOrdersInput{
		UserID:   req.UserId,
		Action:   pvz.ActionType(req.Action),
		OrderIDs: req.OrderIds,
	}
}
func mapProcessResultToProto(res *app.ProcessOrdersOutput) *pkg.ProcessResult {
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

func mapDomainHistoryToProto(his app.OrderHistory) *pkg.OrderHistory {
	return &pkg.OrderHistory{
		OrderId:   his.OrderID,
		Status:    pkg.OrderStatus(his.Status),
		CreatedAt: timestamppb.New(his.CreatedAt),
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
