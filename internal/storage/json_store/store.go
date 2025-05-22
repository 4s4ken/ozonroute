package json_store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/storage/errors"
)

type Storage struct {
	path      string
	store     map[pvz.UserID]map[pvz.OrderID]*pvz.OrderInfo
	history   []OrderHistoryEvent
	orderList map[pvz.UserID]*pvz.OrderInfo
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

func (s *Storage) accept(orderID pvz.OrderID, userID pvz.UserID, expires time.Time) (err error) {
	if expires.Before(time.Now().UTC()) {
		return storageErrors.VALIDATION_FAILED
	}

	if _, ok := s.store[userID]; !ok {
		s.store[userID] = make(map[pvz.OrderID]*pvz.OrderInfo)
	}

	if _, exists := s.store[userID][orderID]; exists {
		return storageErrors.ORDER_ALREADY_EXISTS
	}

	s.store[userID][orderID] = &pvz.OrderInfo{
		Userid:        userID,
		Orderid:       orderID,
		Time:          expires,
		Status:        1,
		StatusChanged: true,
	}
	return s.writeDataToFile()
}

func (s *Storage) Accept(orderID uint64, userID uint64, expires time.Time) error {
	return s.accept(pvz.OrderID(orderID), pvz.UserID(userID), expires)
}

func (s *Storage) returnOrder(orderID pvz.OrderID) error {
	for _, orders := range s.store {
		if order, ok := orders[pvz.OrderID(orderID)]; ok {
			if order.Status == 2 {
				return fmt.Errorf("Invalid order status %d: %w", orderID, storageErrors.VALIDATION_FAILED)

			}
			if order.Status == 1 && order.Time.After(time.Now().UTC()) {
				return fmt.Errorf("order %d storage has expired: %v\n", orderID, storageErrors.STORAGE_EXPIRED)
			}

			order.Status = 3
			now := time.Now().UTC()
			order.ReturnedAt = &now
			order.StatusChanged = true
			return s.writeDataToFile()
		}
	}
	return storageErrors.ORDER_NOT_FOUND
}

func (s *Storage) Return(orderID uint64) error {
	return s.returnOrder(pvz.OrderID(orderID))
}

func (s *Storage) process(userID pvz.UserID, act string, orderIDs []pvz.OrderID) error {
	orders, ok := s.store[userID]
	if !ok {
		return storageErrors.ORDER_NOT_FOUND
	}

	for _, orderID := range orderIDs {
		order, ok := orders[orderID]
		if !ok {
			return storageErrors.ORDER_NOT_FOUND
		}

		switch act {
		case "issue":
			if order.IsExpired(time.Now().UTC()) {
				return fmt.Errorf("order %d storage has expired: %w", orderID, storageErrors.STORAGE_EXPIRED)
			}
			order.Status = 2
			order.StatusChanged = true

		case "return":
			if order.IsExpiredForReturn(time.Now().UTC()) {
				return fmt.Errorf("order %d storage time has expired: %w", orderID, storageErrors.STORAGE_EXPIRED)
			}
			order.Status = 3
			order.StatusChanged = true

		default:
			return storageErrors.VALIDATION_FAILED
		}
	}
	return s.writeDataToFile()
}

func (s *Storage) Process(userID uint64, act string, orderIDs []pvz.OrderID) error {
	pvzOrderIDs := make([]pvz.OrderID, len(orderIDs))
	for i, id := range orderIDs {
		pvzOrderIDs[i] = pvz.OrderID(id)
	}
	return s.process(pvz.UserID(userID), act, pvzOrderIDs)
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
	returned := make([]*pvz.OrderInfo, 0)

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

func (s *Storage) OrderHistory() ([]*pvz.OrderInfo, error) {
	changedOrders := make([]*pvz.OrderInfo, 0)

	for _, userOrders := range s.store {
		for _, order := range userOrders {
			if order.StatusChanged {
				changedOrders = append(changedOrders, order)
			}
		}
	}
	return changedOrders, nil
}

func (s *Storage) ImportOrders(filePath string) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("cannot read file: %w", err)
	}

	temp := make(map[string]map[string]*pvz.OrderInfo)

	if err := json.Unmarshal(data, &temp); err != nil {
		return 0, fmt.Errorf("invalid JSON: %v", storageErrors.VALIDATION_FAILED)
	}

	count := 0
	for userStr, ordersMap := range temp {
		uid, err := strconv.ParseUint(userStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid user id key: %v", storageErrors.VALIDATION_FAILED)
		}
		userID := pvz.UserID(uid)

		if _, exists := s.store[userID]; !exists {
			s.store[userID] = make(map[pvz.OrderID]*pvz.OrderInfo)
		}

		for orderStr, order := range ordersMap {
			oid, err := strconv.ParseUint(orderStr, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid order id key: %v", storageErrors.VALIDATION_FAILED)
			}
			orderID := pvz.OrderID(oid)
			s.store[userID][orderID] = order
			count++
		}
	}
	if err := s.writeDataToFile(); err != nil {
		return 0, storageErrors.INTERNAL_ERROR
	}

	return count, nil
}

func (s *Storage) ScrollOrders(userID uint64, limit uint, lastID *pvz.OrderID) ([]*pvz.OrderInfo, *pvz.OrderID, error) {
	ordersMap, ok := s.store[pvz.UserID(userID)]
	if !ok || len(ordersMap) == 0 {
		return nil, nil, storageErrors.ORDER_NOT_FOUND
	}

	scrollListOrders := make([]*pvz.OrderInfo, 0, len(ordersMap))
	for _, order := range ordersMap {
		scrollListOrders = append(scrollListOrders, order)
	}

	sort.Slice(scrollListOrders, func(i, j int) bool {
		return scrollListOrders[i].Orderid < scrollListOrders[j].Orderid
	})

	startIndex := 0
	if lastID != nil {
		found := false
		for i, order := range scrollListOrders {
			if order.Orderid == *lastID {
				startIndex = i + 1
				found = true
				break
			}
		}
		if !found {
			return nil, nil, fmt.Errorf("lastID %v not found", *lastID)
		}
	}

	endIndex := startIndex + int(limit)
	if endIndex > len(scrollListOrders) {
		endIndex = len(scrollListOrders)
	}
	result := scrollListOrders[startIndex:endIndex]

	var nextLastID *pvz.OrderID
	if endIndex < len(scrollListOrders) {
		nextLastID = &result[len(result)-1].Orderid
	}

	return result, nextLastID, nil
}

func (s *Storage) readDataFromFile(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		s.store = make(map[pvz.UserID]map[pvz.OrderID]*pvz.OrderInfo)
		return nil
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s.store)
	if err != nil {
		if err == io.EOF {
			s.store = make(map[pvz.UserID]map[pvz.OrderID]*pvz.OrderInfo)
			return nil
		}
		return storageErrors.VALIDATION_FAILED
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
		return storageErrors.VALIDATION_FAILED
	}

	return nil
}
