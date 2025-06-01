package pvz

import (
	"time"
)

type (
	OrderID            uint64
	UserID             uint64
	OrderStatus        int
	OrderStatusChanges bool
)

const (
	StatusAccepted OrderStatus = iota
	StatusIssued
	StatusReturned
	StatusExpired
)

func (s OrderStatus) String() string {
	switch s {
	case StatusAccepted:
		return "accepted"
	case StatusIssued:
		return "issued"
	case StatusReturned:
		return "returned"
	case StatusExpired:
		return "expired"
	default:
		return "unknown"
	}
}

type ListOrdersOptions struct {
	UserID uint64
	InPvz  bool
	Last   *uint64
	Page   *uint64
	Limit  *uint64
}

type OrderInfo struct {
	Orderid       OrderID            `json:"order_id,"`
	Userid        UserID             `json:"user_id"`
	Time          time.Time          `json:"date"`
	Status        OrderStatus        `json:"status"`
	ReturnedAt    time.Time          `json:"returned_at,omitempty"`
	StatusChanged OrderStatusChanges `json:"status_changes"`
}

func (o *OrderInfo) IsExpired(now time.Time) bool {
	return time.Time(o.Time).Before(now)
}

func (o *OrderInfo) IsExpiredForReturn(now time.Time) bool {
	twoDaysAfterExpiry := o.Time.Add(48 * time.Hour)
	return now.After(twoDaysAfterExpiry)
}
