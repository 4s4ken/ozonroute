package worker_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
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

func TestOutboxWorker_FirstAndSecondPass_Success(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	payload := []byte(`{"event_id":"ok"}`)

	repo := mockoutbox.NewOutboxRepositoryMock(t)
	producer := &fakeProducer{}

	repo.MsgCreatedMock.Expect(ctx, 10).Return([]outbox.ObEvent{
		{ID: id, Payload: payload},
	}, nil)

	repo.MsgProcessingMock.Expect(ctx, 10).Return([]outbox.ObEvent{
		{ID: id, Payload: payload},
	}, nil)

	var StatusUpdated []struct {
		ID     uuid.UUID
		Status string
	}

	repo.UpdateStatusMock.Set(func(_ context.Context, eid uuid.UUID, status string, _ error) error {
		StatusUpdated = append(StatusUpdated, struct {
			ID     uuid.UUID
			Status string
		}{eid, status})
		return nil
	})

	w := worker.NewOutboxWorker(repo, producer, 10)
	w.TestRun(ctx)

	require.Len(t, producer.Sent, 1)
	require.Equal(t, payload, producer.Sent[0])

	require.Len(t, StatusUpdated, 2)
	require.Equal(t, "PROCESSING", StatusUpdated[0].Status)
	require.Equal(t, "COMPLETED", StatusUpdated[1].Status)
}

func TestOutboxWorker_SecondPass_FailedSend(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	payload := []byte(`{"event_id":"fail"}`)

	repo := mockoutbox.NewOutboxRepositoryMock(t)
	producer := &fakeProducer{Err: errors.New("kafka is down")}

	repo.MsgCreatedMock.Expect(ctx, 10).Return([]outbox.ObEvent{
		{ID: id, Payload: payload},
	}, nil)

	repo.MsgProcessingMock.Expect(ctx, 10).Return([]outbox.ObEvent{
		{ID: id, Payload: payload},
	}, nil)

	var failedStatusUpdated bool
	repo.UpdateStatusMock.Set(func(_ context.Context, eid uuid.UUID, status string, err error) error {
		if status == "FAILED" {
			require.ErrorContains(t, err, "kafka is down")
			failedStatusUpdated = true
		}
		return nil
	})

	w := worker.NewOutboxWorker(repo, producer, 10)
	w.TestRun(ctx)

	require.True(t, failedStatusUpdated)
	require.Len(t, producer.Sent, 0)
}

func TestOutboxWorker_SecondPass_ExceededMaxAttempts(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	payload := []byte(`{"event_id":"retry-3"}`)
	lastAttempt := time.Now().Add(-300 * time.Millisecond)
	event := outbox.ObEvent{
		ID:            id,
		Payload:       payload,
		Status:        pvz.OrderStatusAccepted.String(),
		Attempts:      3,
		LastAttempted: &lastAttempt,
	}

	repo := mockoutbox.NewOutboxRepositoryMock(t)
	producer := &fakeProducer{Err: errors.New("still broken")}

	repo.MsgCreatedMock.Expect(ctx, 10).Return([]outbox.ObEvent{event}, nil)
	repo.MsgProcessingMock.Expect(ctx, 10).Return([]outbox.ObEvent{event}, nil)

	var failed bool

	repo.UpdateStatusMock.Set(func(_ context.Context, eid uuid.UUID, status string, err error) error {
		if status == "FAILED" && err != nil {
			require.EqualError(t, err, "NO_ATTEMPTS_LEFT")
			failed = true
		}
		return nil
	})

	w := worker.NewOutboxWorker(repo, producer, 10)
	w.TestRun(ctx)

	require.True(t, failed)
	require.Len(t, producer.Sent, 0)
}
