package pvz

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"time"
)

func (s *OrdersService) AcceptOrder(ctx context.Context, req *pvz.AcceptOrderRequest) (*pvz.OrderResponse, error) {
	if req == nil {
		return nil, businesserror.EntityValidationError("AcceptOrderRequest", "nil")
	}

	now := time.Now()
	if req.ExpiresAt == nil || req.ExpiresAt.AsTime().Before(now) {
		return nil, businesserror.EntityExpiredError("Order", fmt.Sprintf("%d", req.OrderId))
	}

	if _, exists := s.storage[req.OrderId]; exists {
		return nil, businesserror.EntityAlreadyExistsError("Order", fmt.Sprintf("%d", req.OrderId))
	}

	if req.Package == nil {
		return nil, businesserror.EntityPackageError("Order", fmt.Sprintf("%d", req.OrderId))
	}
	if !isValidWeight(req.Weight, *req.Package) {
		return nil, businesserror.EntityWeightError(string(*req.Package), fmt.Sprintf("%d", req.OrderId))
	}
	totalPrice := req.Price + float32(finalPrice(*req.Package))

	order := &pvz.Order{
		OrderId:    req.OrderId,
		UserId:     req.UserId,
		Status:     pvz.OrderStatus_ORDER_STATUS_ACCEPTED,
		ExpiresAt:  req.ExpiresAt,
		Weight:     req.Weight,
		Package:    *req.Package,
		TotalPrice: totalPrice,
	}

	s.storage[req.OrderId] = order

	return &pvz.OrderResponse{
		Status:  pvz.OrderStatus_ORDER_STATUS_ACCEPTED,
		OrderId: req.OrderId,
	}, nil
}

func uintToStr(id uint64) string {
	return fmt.Sprintf("%d", id)
}

func isValidWeight(weight float32, pkg pvz.PackageType) bool {
	switch pkg {
	case pvz.PackageType_PACKAGE_TYPE_UNSPECIFIED:
		return weight >= 0
	case pvz.PackageType_PACKAGE_TYPE_BAG:
		return weight <= 10
	case pvz.PackageType_PACKAGE_TYPE_BOX:
		return weight <= 30
	case pvz.PackageType_PACKAGE_TYPE_TAPE:
		return weight >= 0
	case pvz.PackageType_PACKAGE_TYPE_BAG_TAPE:
		return weight <= 15
	case pvz.PackageType_PACKAGE_TYPE_BOX_TAPE:
		return weight <= 25
	default:
		return false
	}
}

func finalPrice(pkg pvz.PackageType) int32 {
	switch pkg {
	case pvz.PackageType_PACKAGE_TYPE_BAG:
		return 5
	case pvz.PackageType_PACKAGE_TYPE_BOX:
		return 20
	case pvz.PackageType_PACKAGE_TYPE_TAPE:
		return 1
	case pvz.PackageType_PACKAGE_TYPE_BAG_TAPE:
		return 6
	case pvz.PackageType_PACKAGE_TYPE_BOX_TAPE:

		return 21
	default:
		return 0
	}
}
