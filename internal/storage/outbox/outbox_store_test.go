package storage_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/kafka/worker"
	mockoutbox "gitlab.ozon.dev/gleb_evseev/homework/internal/mocks"
)

type fakeProducer struct {
	Sent [][]byte
	Err  error
}

func (f *fakeProducer) Send(ctx context.Context, key string, value []byte) error {
	if f.Err != nil {
		return f.Err
	}
	f.Sent = append(f.Sent, value)
	return nil
}

func (f *fakeProducer) Close() error {
	return nil
}

func TestOutboxWorker_processBatch_success(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	payload := []byte(`{"event_id": "test"}`)

	repo := mockoutbox.NewOutboxRepositoryMock(t)
	repo.MsgCreatedMock.Return([]outbox.ObEvent{{ID: id, Payload: payload}}, nil)
	repo.MsgProcessingMock.Return([]outbox.ObEvent{{ID: id, Payload: payload}}, nil)
	repo.UpdateStatusMock.Return(nil)

	producer := &fakeProducer{}
	outboxWorker := worker.NewOutboxWorker(repo, producer, 10)

	outboxWorker.TestRun(ctx)

	require.Len(t, producer.Sent, 1)
	require.Equal(t, payload, producer.Sent[0])
}

func TestOutboxWorker_processBatch_failedSend(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	payload := []byte(`{"event_id": "fail"}`)

	repo := mockoutbox.NewOutboxRepositoryMock(t)
	repo.MsgCreatedMock.Return([]outbox.ObEvent{{ID: id, Payload: payload}}, nil)
	repo.MsgProcessingMock.Return([]outbox.ObEvent{{ID: id, Payload: payload}}, nil)
	repo.UpdateStatusMock.Return(nil)

	producer := &fakeProducer{Err: errors.New("kafka down")}
	outboxWorker := worker.NewOutboxWorker(repo, producer, 10)

	outboxWorker.TestRun(ctx)

	require.Len(t, producer.Sent, 0)
}
