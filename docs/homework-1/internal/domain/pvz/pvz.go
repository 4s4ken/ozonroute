package pvz

import (
	"time"
)

type (
	OrderID     uint64
	UserID      uint64
	Date        time.Time
	OrderStatus int
)

const (
	StatusUnknown OrderStatus = iota
	StatusAccepted
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
	ORDERID    OrderID     `json:"order_id,"`
	USERID     UserID      `json:"user_id"`
	Expires    time.Time   `json:"expires_at"`
	Status     OrderStatus `json:"status"`
	ReturnedAt *time.Time  `json:"returned_at"`
}

func (o *OrderInfo) IsExpired(now time.Time) bool {
	return time.Time(o.Expires).Before(now)
}

func (o *OrderInfo) IsExpiredForReturn(now time.Time) bool {
	twoDaysAfterExpiry := o.Expires.Add(-48 * time.Hour)
	return now.After(twoDaysAfterExpiry)
}
