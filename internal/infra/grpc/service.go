package grpc

import (
	"context"
	"errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type OrdersService struct {
	storage map[uint64]*pvz.Order
}

func NewOrdersService() *OrdersService {
	return &OrdersService{
		storage: make(map[uint64]*pvz.Order),
	}
}

func (s *OrdersService) Accept(ctx context.Context, order pvz.AcceptOrderInput) error {

	if _, exists := s.storage[order.OrderID]; exists {
		return errors.New("order already exists")
	}

	pkgType := pvz.PackageTypeUnspecified
	if order.Package != nil {
		pkgType = *order.Package
	}

	s.storage[order.OrderID] = &pvz.Order{
		OrderID:    pvz.OrderID(order.OrderID),
		UserID:     pvz.UserID(order.UserID),
		Status:     pvz.OrderStatusAccepted,
		ExpiresAt:  order.ExpiresAt,
		Weight:     order.Weight,
		TotalPrice: order.Price,
		Package:    pkgType,
	}

	return nil
}

func (s *OrdersService) Return(ctx context.Context, orderID pvz.ReturnOrderInput) error {
	order, exists := s.storage[uint64(orderID.OrderID)]
	if !exists {
		return errors.New("order not found")
	}

	order.Status = pvz.OrderStatusReturned
	return nil
}

func (s *OrdersService) Process(ctx context.Context, orders pvz.ProcessOrdersInput) (pvz.ProcessResult, error) {

	var processed []uint64
	var errorsList []uint64

	for _, id := range orders.OrderIDs {
		order, exists := s.storage[uint64(id)]
		if !exists {
			errorsList = append(errorsList, uint64(id))
			continue
		}

		switch orders.Action {
		case pvz.ActionTypeIssue:
			order.Status = pvz.OrderStatusIssued
			processed = append(processed, uint64(id))
		case pvz.ActionTypeReturn:
			order.Status = pvz.OrderStatusReturned
			processed = append(processed, uint64(id))
		default:
			errorsList = append(errorsList, uint64(id))
		}
	}

	return pvz.ProcessResult{
		Processed: processed,
		Errors:    errorsList,
	}, nil
}

func (s *OrdersService) ListOrders(ctx context.Context, opts pvz.ListOrdersInput) (*pvz.OrdersList, error) {
	orders := make([]pvz.Order, 0)
	for _, order := range s.storage {
		if order.UserID == opts.UserID {
			orders = append(orders, *order)
		}
	}

	total := int32(len(orders))

	if opts.Pagination != nil {
		start := int(opts.Pagination.Page * opts.Pagination.CountOnPage)
		end := start + int(opts.Pagination.CountOnPage)
		if start > len(orders) {
			orders = []pvz.Order{}
		} else if end > len(orders) {
			orders = orders[start:]
		} else {
			orders = orders[start:end]
		}
	}

	return &pvz.OrdersList{
		Orders: orders,
		Total:  total,
	}, nil
}

func (s *OrdersService) ListReturns(ctx context.Context, opts pvz.ListReturnsInput) (*pvz.ReturnsList, error) {
	returned := make([]pvz.Order, 0)
	for _, order := range s.storage {
		if order.Status == pvz.OrderStatusReturned {
			returned = append(returned, *order)
		}
	}

	return &pvz.ReturnsList{
		Returns: returned,
	}, nil
}

func (s *OrdersService) OrderHistory(ctx context.Context) (*pvz.OrderHistoryList, error) {
	return &pvz.OrderHistoryList{
		History: []pvz.OrderHistory{},
	}, nil
}
