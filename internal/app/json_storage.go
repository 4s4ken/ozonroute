package app

import (
	"encoding/json"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// Storage представляет хранилище заказов, сохраняемое в JSON-файле.
type Storage struct {
	file  *os.File
	path  string
	store map[pvz.UserID]map[pvz.OrderID]*pvz.Order
	nowFn func() time.Time
}

// NewJsonStorage создаёт новое хранилище заказов, связанное с указанным JSON-файлом.
func NewJsonStorage(path string, nowFn func() time.Time) (*Storage, error) {

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("create directory for db file: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		path:  path,
		store: make(map[pvz.UserID]map[pvz.OrderID]*pvz.Order),
		file:  f,
		nowFn: nowFn,
	}
	err = s.readDataFromFile(f)
	if err != nil {
		return nil, storageErrors.EntityInternalError("Json file", "fail to open json storage")
	}
	return s, nil
}

// readDataFromFile загружает данные из JSON-файла в память.
func (s *Storage) readDataFromFile(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		s.store = make(map[pvz.UserID]map[pvz.OrderID]*pvz.Order)
		return nil
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s.store)
	if err != nil {
		if err == io.EOF {
			s.store = make(map[pvz.UserID]map[pvz.OrderID]*pvz.Order)
			return nil
		}
		return storageErrors.EntityValidationError("JSON file", "fail reading")
	}
	return nil
}

// writeDataToFile сериализует текущее состояние in-memory хранилища (s.store)
func (s *Storage) writeDataToFile() error {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {

			fmt.Printf("failed to close file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s.store)
	if err != nil {
		return storageErrors.EntityInternalError("JSON file", "fail writing")
	}

	return nil
}

func (s *Storage) Close() error {
	return s.file.Close()
}
