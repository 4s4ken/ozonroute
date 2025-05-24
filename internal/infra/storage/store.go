package storage

import (
	"encoding/json"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"io"
	"os"
	"time"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type Storage struct {
	file      *os.File
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
		_ = f.Close()
	}(f)

	s := &Storage{
		path:  path,
		store: make(map[pvz.UserID]map[pvz.OrderID]*pvz.OrderInfo),
		file:  f,
	}
	err = s.readDataFromFile(f)
	if err != nil {
		return nil, err
	}

	return s, nil
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
		return storageErrors.EntityValidationError("JSON file", "fail reading")
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
		_ = file.Close()
	}(file)

	encoder := json.NewEncoder(file)

	err = encoder.Encode(s.store)
	if err != nil {
		return storageErrors.EntityValidationError("JSON file", "fail writing")
	}

	return nil
}

func (s *Storage) Close() error {
	return s.file.Close()
}
