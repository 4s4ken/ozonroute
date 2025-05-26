package usecase

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"time"
)

type storage interface {
	Accept(ctx context.Context, orderID pvz.OrderID, userID pvz.UserID, exp time.Time) error
	Return(ctx context.Context, orderID pvz.OrderID) error
	Process(ctx context.Context, userID pvz.UserID, act string, orderIDs []pvz.OrderID) error
	ListOrders(ctx context.Context, typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	ListReturns(ctx context.Context, typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	OrderHistory() ([]*pvz.OrderInfo, error)
	ImportOrders(ctx context.Context, orders string) (int, error)
	ScrollOrders(ctx context.Context, userID pvz.UserID, limit uint, lastID *pvz.OrderID) ([]*pvz.OrderInfo, *pvz.OrderID, error)
}

type Service struct {
	storage storage
}

func New(storage storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Accept(ctx context.Context, orderID pvz.OrderID, userID pvz.UserID, exp time.Time) error {

	err := s.storage.Accept(ctx, orderID, userID, exp)
	if err != nil {
		return fmt.Errorf("storage.Accept: %w", err)
	}

	return nil
}

func (s *Service) Return(ctx context.Context, orderID pvz.OrderID) error {
	err := s.storage.Return(ctx, orderID)
	if err != nil {
		return fmt.Errorf("storage.Return: %w", err)
	}

	return nil
}

func (s *Service) Process(ctx context.Context, userID pvz.UserID, act string, orderIDs []pvz.OrderID) error {

	err := s.storage.Process(ctx, userID, act, orderIDs)
	if err != nil {
		return fmt.Errorf("storage.Process: %w", err)
	}
	return nil
}

func (s *Service) ListOrders(ctx context.Context, opts pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {
	orders, err := s.storage.ListOrders(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.ListOrders: %w", err)
	}
	return orders, nil
}

func (s *Service) ListReturns(ctx context.Context, opts pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {
	returns, err := s.storage.ListReturns(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.ListReturns: %w", err)
	}
	return returns, nil
}

func (s *Service) OrderHistory() ([]*pvz.OrderInfo, error) {
	history, err := s.storage.OrderHistory()
	if err != nil {
		return nil, fmt.Errorf("storage.OrderHistory: %w", err)
	}
	return history, nil
}

func (s *Service) ImportOrders(ctx context.Context, orders string) (int, error) {

	count, err := s.storage.ImportOrders(ctx, orders)
	if err != nil {
		return 0, fmt.Errorf("storage.ImportOrders: %w", err)
	}
	return count, nil
}

func (s *Service) ScrollOrders(ctx context.Context, userID pvz.UserID, limit uint, lastID *pvz.OrderID) ([]*pvz.OrderInfo, *pvz.OrderID, error) {
	orders, nextID, err := s.storage.ScrollOrders(ctx, userID, limit, lastID)
	if err != nil {
		return nil, nil, fmt.Errorf("storage.ScrollOrders: %w", err)
	}
	return orders, nextID, nil
}
