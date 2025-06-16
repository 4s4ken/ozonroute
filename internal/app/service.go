package app

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type storage interface {
	Accept(ctx context.Context, order pvz.OrderCreate) error
	Return(ctx context.Context, input pvz.OrderCreate) error
	Process(ctx context.Context, orders ProcessOrdersInput) (*ProcessOrdersOutput, error)
	ListOrders(ctx context.Context, opts ListOrdersInput) (*OrdersListOutput, error)
	ListReturns(ctx context.Context, opts ListReturnsInput) (*ReturnsListOutput, error)
	OrderHistory(ctx context.Context) (*OrderHistoryList, error)
}

func New(storage *Storage) *Service {
	return &Service{
		storage: storage,
	}
}

type Service struct {
	storage storage
}

func (s *Service) Accept(ctx context.Context, order pvz.OrderCreate) error {
	err := s.storage.Accept(ctx, order)
	if err != nil {
		return fmt.Errorf("storage.Accept: %w", err)
	}
	return nil
}

func (s *Service) Return(ctx context.Context, orderID pvz.OrderCreate) error {
	err := s.storage.Return(ctx, orderID)
	if err != nil {
		return fmt.Errorf("storage.Return: %w", err)
	}
	return nil
}

func (s *Service) Process(ctx context.Context, orders ProcessOrdersInput) (*ProcessOrdersOutput, error) {
	ordersID, err := s.storage.Process(ctx, orders)
	if err != nil {
		return ordersID, fmt.Errorf("storage.Process: %w", err)
	}
	return ordersID, nil
}

func (s *Service) ListOrders(ctx context.Context, opts ListOrdersInput) (*OrdersListOutput, error) {
	orders, err := s.storage.ListOrders(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.ListOrders: %w", err)
	}
	return orders, nil
}

func (s *Service) ListReturns(ctx context.Context, opts ListReturnsInput) (*ReturnsListOutput, error) {
	returns, err := s.storage.ListReturns(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.ListReturns: %w", err)
	}
	return returns, nil
}

func (s *Service) OrderHistory(ctx context.Context) (*OrderHistoryList, error) {
	history, err := s.storage.OrderHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.OrderHistory: %w", err)
	}
	return history, nil
}
