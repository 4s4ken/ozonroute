package app

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type storage interface {
	Accept(ctx context.Context, order pvz.AcceptOrderInput) error
	Return(ctx context.Context, orderID pvz.ReturnOrderInput) error
	Process(ctx context.Context, orders pvz.ProcessOrdersInput) (pvz.ProcessResult, error)
	ListOrders(ctx context.Context, opts pvz.ListOrdersInput) (*pvz.OrdersList, error)
	ListReturns(ctx context.Context, opts pvz.ListReturnsInput) (*pvz.ReturnsList, error)
	OrderHistory(ctx context.Context) (*pvz.OrderHistoryList, error)
	//ImportOrders(ctx context.Context, orders string) (pvz.ImportResult, error)
	//ScrollOrders(ctx context.Context, userID pvz.UserID, limit uint, lastID *pvz.OrderID) ([]*pvz.OrdersList, *pvz.OrderID, error)
}

func New(storage storage) *Service {
	return &Service{
		storage: storage,
	}
}

type Service struct {
	storage storage
}

func (s *Service) Accept(ctx context.Context, order pvz.AcceptOrderInput) error {
	err := s.storage.Accept(ctx, order)
	if err != nil {
		return fmt.Errorf("storage.Accept: %w", err)
	}
	return nil
}

func (s *Service) Return(ctx context.Context, orderID pvz.ReturnOrderInput) error {
	err := s.storage.Return(ctx, orderID)
	if err != nil {
		return fmt.Errorf("storage.Return: %w", err)
	}
	return nil
}

func (s *Service) Process(ctx context.Context, orders pvz.ProcessOrdersInput) (pvz.ProcessResult, error) {
	ordersID, err := s.storage.Process(ctx, orders)
	if err != nil {
		return ordersID, fmt.Errorf("storage.Process: %w", err)
	}
	return ordersID, nil
}

func (s *Service) ListOrders(ctx context.Context, opts pvz.ListOrdersInput) (*pvz.OrdersList, error) {
	orders, err := s.storage.ListOrders(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.ListOrders: %w", err)
	}
	return orders, nil
}

func (s *Service) ListReturns(ctx context.Context, opts pvz.ListReturnsInput) (*pvz.ReturnsList, error) {
	returns, err := s.storage.ListReturns(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.ListReturns: %w", err)
	}
	return returns, nil
}

func (s *Service) OrderHistory(ctx context.Context) (*pvz.OrderHistoryList, error) {
	history, err := s.storage.OrderHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.OrderHistory: %w", err)
	}
	return history, nil
}
