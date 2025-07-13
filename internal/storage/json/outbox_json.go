package json

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"os"
	"path/filepath"
	"sync"
)

type storage struct {
	mu       sync.Mutex
	filepath string
}

func New(path string) outbox.OutboxRepository {
	return &storage{filepath: path}
}

func (s *storage) SaveEvent(ctx context.Context, e outbox.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	events, _ := s.read()
	events = append(events, e)
	return s.write(events)
}

func (s *storage) read() ([]outbox.Event, error) {
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return []outbox.Event{}, nil
		}
		return nil, err
	}

	var events []outbox.Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (s *storage) write(events []outbox.Event) error {
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(s.filepath), os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filepath, data, 0644)
}

func (s *storage) MsgCreated(ctx context.Context, limit int) ([]outbox.ObEvent, error) {
	return nil, nil
}

func (s *storage) MsgProcessing(ctx context.Context, limit int) ([]outbox.ObEvent, error) {
	return nil, nil
}

func (s *storage) UpdateStatus(ctx context.Context, id uuid.UUID, status string, reason error) error {
	return nil
}
