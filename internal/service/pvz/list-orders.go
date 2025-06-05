package pvz

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

func (s *OrdersService) ListOrders(ctx context.Context, req *pvz.ListOrdersRequest) (*pvz.OrdersList, error) {
	if req == nil || req.UserId == 0 {
		return nil, businesserror.EntityValidationError("ListOrdersRequest", "")
	}

	userOrders := make([]*pvz.Order, 0)
	for _, order := range s.storage {
		if order.UserId == req.UserId {
			if req.InPvz {
				if order.Status == pvz.OrderStatus_ORDER_STATUS_ACCEPTED {
					userOrders = append(userOrders, order)
				}
			} else {
				userOrders = append(userOrders, order)
			}
		}
	}

	if req.LastN != nil && *req.LastN > 0 && int(*req.LastN) < len(userOrders) {
		userOrders = userOrders[len(userOrders)-int(*req.LastN):]
	}

	if req.Pagination == nil {
		return &pvz.OrdersList{
			Orders: userOrders,
			Total:  int32(len(userOrders)),
		}, nil
	}

	page := uint32(1)
	countOnPage := uint32(len(userOrders))

	if req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}
	if req.Pagination.CountOnPage > 0 {
		countOnPage = req.Pagination.CountOnPage
	}

	start := (page - 1) * countOnPage
	if int(start) > len(userOrders) {
		start = uint32(len(userOrders))
	}
	end := start + countOnPage
	if int(end) > len(userOrders) {
		end = uint32(len(userOrders))
	}

	ordersPage := userOrders[start:end]

	return &pvz.OrdersList{
		Orders: ordersPage,
		Total:  int32(len(userOrders)),
	}, nil
}
