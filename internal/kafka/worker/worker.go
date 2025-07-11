package worker

import (
	"context"
	"errors"
	prod "gitlab.ozon.dev/gleb_evseev/homework/internal/kafka/producer"
	"log"
	"time"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
)

type OutboxWorker struct {
	outboxRepo outbox.OutboxRepository
	producer   prod.Producer
	batchSize  int
}

func NewOutboxWorker(repo outbox.OutboxRepository, producer prod.Producer, batchSize int) *OutboxWorker {
	return &OutboxWorker{
		outboxRepo: repo,
		producer:   producer,
		batchSize:  batchSize,
	}
}

func (w *OutboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Outbox worker stopped")
			return
		case <-ticker.C:
			w.firstPass(ctx)
			w.secondPass(ctx)
		}
	}
}

func (w *OutboxWorker) firstPass(ctx context.Context) {
	events, err := w.outboxRepo.MsgCreated(ctx, w.batchSize)
	if err != nil {
		log.Printf("[FirstPass] Failed to CREATED events: %v", err)
		return
	}
	if len(events) == 0 {
		return
	}

	for _, event := range events {
		err := w.outboxRepo.UpdateStatus(ctx, event.ID, "PROCESSING", nil)
		if err != nil {
			log.Printf("[FirstPass] Failed to PROCESSING event %s: %v", event.ID, err)
		}
	}
}

func (w *OutboxWorker) secondPass(ctx context.Context) {
	events, err := w.outboxRepo.MsgProcessing(ctx, w.batchSize)
	if err != nil {
		log.Printf("[SecondPass] Failed to PROCESSING events: %v", err)
		return
	}
	if len(events) == 0 {
		return
	}

	for _, event := range events {
		if event.Attempts >= 3 {
			log.Printf("[SecondPass] Event %s exceeded max attempts", event.ID)
			_ = w.outboxRepo.UpdateStatus(ctx, event.ID, "FAILED", errors.New("NO_ATTEMPTS_LEFT"))
			continue
		}

		err := w.producer.Send(ctx, event.ID.String(), event.Payload)
		if err != nil {
			log.Printf("[SecondPass] Failed to send event %s: %v", event.ID, err)
			_ = w.outboxRepo.UpdateStatus(ctx, event.ID, "FAILED", err)
			continue
		}

		err = w.outboxRepo.UpdateStatus(ctx, event.ID, "COMPLETED", nil)
		if err != nil {
			log.Printf("[SecondPass] Failed to COMPLETED event %s: %v", event.ID, err)
		}
	}
}

func (w *OutboxWorker) TestRun(ctx context.Context) {
	w.firstPass(ctx)
	w.secondPass(ctx)
}
