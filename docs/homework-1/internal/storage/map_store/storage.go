package map_store

import (
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/domain/pvz"
	storage_errors "gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/storage/errors"
	"go.uber.org/multierr"
	"os"
	"strconv"
	"time"
)

type OrderHistoryEvent struct {
	OrderID pvz.OrderID
	Status  pvz.OrderStatus
	Data    time.Time
}

type ListReturns struct {
	OrderID pvz.OrderID
	UserId  pvz.UserID
	Returnd time.Time
}

type Storage struct {
	store     map[pvz.UserID]map[pvz.OrderID]*pvz.OrderInfo
	orderList map[pvz.UserID][]pvz.OrderID
	history   []OrderHistoryEvent
}

func NewStorage() *Storage {
	return &Storage{
		store:     make(map[pvz.UserID]map[pvz.OrderID]*pvz.OrderInfo),
		orderList: make(map[pvz.UserID][]pvz.OrderID),
	}
}

func (s *Storage) accept(orderID pvz.OrderID, userID pvz.UserID, expires time.Time) error {
	if expires.Before(time.Now()) {
		return storage_errors.STORAGE_EXPIRED
	}

	if _, ok := s.store[userID]; !ok {
		s.store[userID] = make(map[pvz.OrderID]*pvz.OrderInfo)
		s.orderList[userID] = []pvz.OrderID{}
	}

	if _, exists := s.store[userID][orderID]; exists {
		return storage_errors.ORDER_ALREADY_EXISTS
	}
	fmt.Println("ORDER_ACCEPTED:", orderID)
	s.store[userID][orderID] = &pvz.OrderInfo{
		ORDERID: orderID,
		USERID:  userID,
		Expires: expires,
		Status:  pvz.StatusAccepted,
	}
	s.orderList[userID] = append(s.orderList[userID], orderID)

	s.history = append(s.history, OrderHistoryEvent{
		OrderID: orderID,
		Status:  pvz.StatusAccepted,
		Data:    time.Now(),
	})

	return nil
}

func (s *Storage) Accept(orderID pvz.OrderID, userID pvz.UserID, date time.Time) error {
	return s.accept(orderID, userID, date)
}

func (s *Storage) returnOrd(orderID pvz.OrderID) error {
	userID := pvz.UserID(orderID)
	order, ok := s.store[userID][orderID]
	if !ok {
		return storage_errors.ORDER_NOT_FOUND
	}

	if order.Status == pvz.StatusIssued {
		return fmt.Errorf("The order %d has already been issued, %w", orderID, storage_errors.VALIDATION_FAILED)
	}

	order.Status = pvz.StatusReturned
	returnTime := time.Now()
	order.ReturnedAt = &returnTime
	fmt.Println("ORDER_RETURNED:", orderID)
	s.history = append(s.history, OrderHistoryEvent{
		OrderID: orderID,
		Status:  3,
		Data:    returnTime,
	})

	return nil
}

func (s *Storage) Return(orderID pvz.OrderID) error {
	return s.returnOrd(orderID)
}

func (s *Storage) process(userID pvz.UserID, action string, orderIDs []pvz.OrderID) error {
	var resultErr error

	for _, orderID := range orderIDs {
		order, ok := s.store[userID][orderID]
		if !ok {
			resultErr = multierr.Append(resultErr, fmt.Errorf("Order %d not found: %v", orderID, storage_errors.ORDER_NOT_FOUND))
			continue
		}

		switch action {
		case "issue":
			if order.IsExpired(time.Now()) {
				resultErr = multierr.Append(resultErr, fmt.Errorf("Order %d expired: %v", orderID, storage_errors.STORAGE_EXPIRED))
				order.Status = pvz.StatusExpired
				continue
			}
			order.Status = pvz.StatusIssued
			s.history = append(s.history, OrderHistoryEvent{OrderID: orderID, Status: order.Status, Data: time.Now()})

		case "return":
			if order.IsExpiredForReturn(time.Now()) {
				resultErr = multierr.Append(resultErr, fmt.Errorf("Order %d expired for return: %v", orderID, storage_errors.STORAGE_EXPIRED))
				order.Status = pvz.StatusExpired
				continue
			}
			order.Status = pvz.StatusReturned
			s.history = append(s.history, OrderHistoryEvent{OrderID: orderID, Status: order.Status, Data: time.Now()})

		default:
			return storage_errors.VALIDATION_FAILED
		}
	}
	return resultErr
}

func (s *Storage) Process(userID pvz.UserID, action string, orderIDs []pvz.OrderID) error {
	return s.process(userID, action, orderIDs)
}

func (s *Storage) listOrders(options pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {
	ordersMap, ok := s.store[pvz.UserID(options.UserID)]
	if !ok || len(ordersMap) == 0 {
		return nil, nil
	}

	listOrd := make([]*pvz.OrderInfo, 0, len(ordersMap))
	for _, order := range ordersMap {
		if options.InPvz && order.Status != pvz.StatusAccepted {
			continue
		}
		listOrd = append(listOrd, order)
	}

	if options.Page != nil && options.Limit != nil {
		start := int(*options.Page) * int(*options.Limit)
		end := start + int(*options.Limit)
		if start < len(listOrd) {
			if end > len(listOrd) {
				end = len(listOrd)
			}
			listOrd = listOrd[start:end]
		} else {
			listOrd = []*pvz.OrderInfo{}
		}
	}

	if options.Last != nil {
		n := int(*options.Last)
		if n < len(listOrd) {
			listOrd = listOrd[len(listOrd)-n:]
		}
	}

	return listOrd, nil
}

func (s *Storage) ListOrders(options pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {
	return s.listOrders(options)
}

func (s *Storage) ListReturns(options pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {
	return s.listReturns(options)
}

func (s *Storage) listReturns(options pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {
	var returned []*pvz.OrderInfo

	for _, userOrders := range s.store {
		for _, order := range userOrders {
			if order.Status == 3 {
				returned = append(returned, order)
			}
		}
	}

	if options.Page != nil && options.Limit != nil {
		start := int(*options.Page) * int(*options.Limit)
		end := start + int(*options.Limit)

		if start < len(returned) {
			if end > len(returned) {
				end = len(returned)
			}
			return returned[start:end], nil
		}
		return []*pvz.OrderInfo{}, nil
	}

	return returned, nil
}

func (s *Storage) orderHistory() ([]OrderHistoryEvent, error) {
	return s.history, nil
}

func (s *Storage) OrderHistory() ([]OrderHistoryEvent, error) {
	return s.orderHistory()
}

func (s *Storage) ImportOrdersFromFile(filePath string) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("cannot read file: %w", err)
	}
	temp := make(map[string]map[string]*pvz.OrderInfo)

	if err := json.Unmarshal(data, &temp); err != nil {
		return 0, fmt.Errorf("invalid JSON: %w", err)
	}

	count := 0
	for userStr, ordersMap := range temp {
		uid, err := strconv.ParseUint(userStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid user id key: %s", userStr)
		}
		userID := pvz.UserID(uid)

		if _, exists := s.store[userID]; !exists {
			s.store[userID] = make(map[pvz.OrderID]*pvz.OrderInfo)
		}

		for orderStr, order := range ordersMap {
			oid, err := strconv.ParseUint(orderStr, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid order id key: %s", orderStr)
			}
			orderID := pvz.OrderID(oid)
			s.store[userID][orderID] = order
			count++
		}
	}

	return count, nil
}
