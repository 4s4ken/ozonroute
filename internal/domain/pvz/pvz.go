package pvz

import (
	"github.com/shopspring/decimal"
	"time"
)

type (
	PackageType string
	OrderStatus int
	ActionType  int
	UserID      uint64
	OrderID     uint64
	OrderIDS    []OrderID
)

const (
	NoPackage          PackageType = ""
	PackageTypeBag     PackageType = "bag"
	PackageTypeBox     PackageType = "box"
	PackageTypeTape    PackageType = "tape"
	PackageTypeBagTape PackageType = "bag+tape"
	PackageTypeBoxTape PackageType = "box+tape"
)

const (
	OrderStatusUnspecified OrderStatus = iota
	OrderStatusAccepted
	OrderStatusIssued
	OrderStatusReturned
	OrderStatusDeleted
)

const (
	ActionTypeUnspecified ActionType = iota
	ActionTypeIssue
	ActionTypeReturn
)

type OrderCreate struct {
	OrderID   OrderID         `json:"order_id,"`
	UserID    UserID          `json:"user_id"`
	ExpiresAt time.Time       `json:"expires_at"`
	Weight    float32         `json:"weight"`
	Price     decimal.Decimal `json:"price"`
	Package   *PackageType    `json:"package_type,omitempty"`
}

type Order struct {
	OrderCreate
	Status     OrderStatus `json:"status"`
	LastUpdate time.Time   `json:"last_update,omitempty"`
}
type StatusChange struct {
	OrderID OrderID
	Status  OrderStatus
	Time    time.Time
}
type ListOptions struct {
	UserID     uint64
	InPVZ      bool
	LastN      *uint64
	Pagination *Pagination
}
type Pagination struct {
	Page  *uint64
	Limit *uint64
}

// IsExpired проверяет, истёк ли срок хранения заказа.
func (o *Order) IsExpired(now time.Time) bool {
	return time.Time(o.ExpiresAt).Before(now)
}

// IsExpiredForReturn проверяет, истёк ли срок возврата заказа.
func (o *Order) IsExpiredForReturn(now time.Time) bool {
	twoDaysAfterExpiry := o.ExpiresAt.Add(48 * time.Hour)
	return now.After(twoDaysAfterExpiry)
}
