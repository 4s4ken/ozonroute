package pvz

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"time"
)

// ReturnOrder возвращает просроченный заказ курьеру.
func (s *OrdersService) ReturnOrder(ctx context.Context, req *pvz.OrderIdRequest) (*pvz.OrderResponse, error) {
	if req == nil {
		return nil, businesserror.EntityValidationError("OrderIdRequest", "empty")
	}

	order, exists := s.storage[req.OrderId]
	if !exists {
		return nil, businesserror.EntityNotFoundError("Order", fmt.Sprintf("%d", req.OrderId))
	}

	if order.Status == pvz.OrderStatus_ORDER_STATUS_RETURNED {
		return nil, businesserror.EntityAlreadyExistsError("Order", fmt.Sprintf("%d", req.OrderId))
	}

	if order.ExpiresAt != nil && order.ExpiresAt.AsTime().Before(time.Now()) {
		return nil, businesserror.EntityExpiredError("Order", fmt.Sprintf("%d", req.OrderId))
	}

	order.Status = pvz.OrderStatus_ORDER_STATUS_RETURNED
	s.storage[req.OrderId] = order

	return &pvz.OrderResponse{
		Status:  pvz.OrderStatus_ORDER_STATUS_RETURNED,
		OrderId: req.OrderId,
	}, nil
}
