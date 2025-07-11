package outbox

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type (
	EventType string
	ActorType string
)

const (
	OrderAccepted          EventType = "order_accepted"
	OrderReturnedToCourier EventType = "order_returned_to_courier"
	OrderIssued            EventType = "order_issued"
	OrderReturnedByClient  EventType = "order_returned_by_client"
)

const (
	ActorCourier ActorType = "courier"
	ActorClient  ActorType = "client"
)

type Event struct {
	EventID   uuid.UUID `json:"event_id"`
	EventType EventType `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Actor     Actor     `json:"actor"`
	Order     Order     `json:"order"`
	Source    string    `json:"source"`
}

type Actor struct {
	Type ActorType `json:"type"`
	ID   int       `json:"id"`
}

type Order struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type ObEvent struct {
	ID            uuid.UUID
	Payload       []byte
	Status        string
	Error         *string
	CreatedAt     time.Time
	SentAt        *time.Time
	Attempts      int
	LastAttempted *time.Time
}

type OutboxRepository interface {
	MsgCreated(ctx context.Context, limit int) ([]ObEvent, error)
	SaveEvent(ctx context.Context, event Event) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, reason error) error
	MsgProcessing(ctx context.Context, limit int) ([]ObEvent, error)
}
