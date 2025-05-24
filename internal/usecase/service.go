package usecase

import (
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"time"
)

type storage interface {
	Accept(orderID uint64, userID uint64, exp time.Time) error
	Return(orderID uint64) error
	Process(userID uint64, act string, orderIDs []pvz.OrderID) error
	ListOrders(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	ListReturns(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	OrderHistory() ([]*pvz.OrderInfo, error)
	ImportOrders(orders string) (int, error)
	ScrollOrders(userID uint64, limit uint, lastID *pvz.OrderID) ([]*pvz.OrderInfo, *pvz.OrderID, error)
}

type Service struct {
	storage storage
}

func New(storage storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Accept(orderID uint64, userID uint64, exp time.Time) error {
	//  валидации

	err := s.storage.Accept(orderID, userID, exp)
	if err != nil {
		return fmt.Errorf("storage.Accept: %w", err)
	}

	return nil
}

func (s *Service) Return(orderID uint64) error {

}
func (s *Service) Process(userID uint64, act string, orderIDs []pvz.OrderID) error {

}
func (s *Service) ListOrders(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {

}
func (s *Service) ListReturns(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {

}
func (s *Service) OrderHistory() ([]*pvz.OrderInfo, error) {

}
func (s *Service) ImportOrders(orders string) (int, error) {

}
func (s *Service) ScrollOrders(userID uint64, limit uint, lastID *pvz.OrderID) ([]*pvz.OrderInfo, *pvz.OrderID, error) {
}
