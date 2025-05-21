package json_store

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/domain/pvz"
	storage_errors "gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/storage/errors"
)

type Storage struct {
	path    string
	store   map[pvz.UserID]map[pvz.OrderID]*pvz.OrderInfo
	history []OrderHistoryEvent
}

type OrderHistoryEvent struct {
	OrderID pvz.OrderID
	Status  pvz.OrderStatus
	Data    time.Time
}

func NewJsonStorage(path string) (*Storage, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	s := &Storage{
		path:  path,
		store: make(map[pvz.UserID]map[pvz.OrderID]*pvz.OrderInfo),
	}
	err = s.readDataFromFile(f)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) Accept(orderID uint64, userID uint64, expires time.Time) error {
	return s.accept(pvz.OrderID(orderID), pvz.UserID(userID), expires)
}

func (s *Storage) accept(orderID pvz.OrderID, userID pvz.UserID, expires time.Time) (err error) {

	if expires.Before(time.Now()) {
		return storage_errors.STORAGE_EXPIRED
	}

	if _, ok := s.store[userID]; !ok {
		s.store[userID] = make(map[pvz.OrderID]*pvz.OrderInfo)
	}

	if _, exists := s.store[userID][orderID]; exists {
		return storage_errors.ORDER_ALREADY_EXISTS
	}
	if _, ok := s.store[userID][orderID]; !ok {
		s.store[userID][orderID] = &pvz.OrderInfo{
			ORDERID: orderID,
			Expires: expires,
			Status:  pvz.StatusAccepted,
		}
	}
	return s.writeDataToFile()
}

func (s *Storage) Return(orderID uint64) error {
	for _, orders := range s.store {
		if order, ok := orders[pvz.OrderID(orderID)]; ok {
			if order.Status == pvz.StatusIssued {
				return storage_errors.VALIDATION_FAILED
			}

			order.Status = pvz.StatusReturned
			now := time.Now()
			order.ReturnedAt = &now

			s.history = append(s.history, OrderHistoryEvent{
				OrderID: pvz.OrderID(orderID),
				Status:  pvz.StatusReturned,
				Data:    now,
			})

			return s.writeDataToFile()
		}
	}
	return storage_errors.ORDER_NOT_FOUND
}

func (s *Storage) Process(userID uint64, act string, orderIDs []pvz.OrderID) error {
	for _, orderID := range orderIDs {
		order, ok := s.store[pvz.UserID(userID)][orderID]
		if !ok {
			fmt.Println(fmt.Errorf("order %d not found: %v", orderID, storage_errors.ORDER_NOT_FOUND))
		}

		switch act {
		case "issue":
			if order.IsExpired(time.Now().UTC()) {
				fmt.Printf("order %d storage time has expired: %v\n", orderID, storage_errors.STORAGE_EXPIRED)
			}
			order.Status = 2
			s.history = append(s.history, OrderHistoryEvent{OrderID: orderID, Status: order.Status, Data: time.Now()})
		case "return":
			if order.IsExpiredForReturn(time.Now()) {
				fmt.Printf("order %d storage time has expired: %v\n", orderID, storage_errors.STORAGE_EXPIRED)
			}
			order.Status = 3
			s.history = append(s.history, OrderHistoryEvent{OrderID: orderID, Status: order.Status, Data: time.Now()})

		default:
			fmt.Printf("invalid action: %v", storage_errors.VALIDATION_FAILED)
		}
	}
	return nil
}

func (s *Storage) ListOrders(options pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {
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

func (s *Storage) ListReturns(options pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {
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

func (s *Storage) OrderHistory() ([]OrderHistoryEvent, error) {
	return s.history, nil
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
	if err := s.writeDataToFile(); err != nil {
		return 0, fmt.Errorf("write error: %w", err)
	}

	return count, nil
}

func (s *Storage) readDataFromFile(file *os.File) (err error) {
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&s.store)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json: %s, %w", s.path, err)
	}

	return nil
}

func (s *Storage) writeDataToFile() (err error) {
	var file *os.File
	file, err = os.OpenFile(s.path, os.O_RDWR|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	encoder := json.NewEncoder(file)

	err = encoder.Encode(s.store)
	if err != nil {
		return fmt.Errorf("cannot write file: %s, %w", s.path, err)
	}

	return nil
}
