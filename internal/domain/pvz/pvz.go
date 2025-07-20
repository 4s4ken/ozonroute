package pvz

import (
	"github.com/shopspring/decimal"
	"time"
)

type (
	PackageType int32
	OrderStatus int
	ActionType  int
	UserID      uint64
	OrderID     uint64
	OrderIDS    []OrderID
)

// PackageType — параметры упаковки заказов.
const (
	NoPackage          PackageType = 0
	PackageTypeBag     PackageType = 1
	PackageTypeBox     PackageType = 2
	PackageTypeTape    PackageType = 3
	PackageTypeBagTape PackageType = 4
	PackageTypeBoxTape PackageType = 5
)

// OrderStatus — параметры статуса заказа.
const (
	OrderStatusUnspecified OrderStatus = iota
	OrderStatusAccepted
	OrderStatusIssued
	OrderStatusReturned
	OrderStatusDeleted
)

// ActionType — параметры статуса обработки заказа.
const (
	ActionTypeUnspecified ActionType = iota
	ActionTypeIssue
	ActionTypeReturn
)

// OrderCreate — параметры заказа при создании.
type OrderCreate struct {
	OrderID   OrderID         `json:"order_id,"`
	UserID    UserID          `json:"user_id"`
	ExpiresAt time.Time       `json:"expires_at"`
	Weight    float32         `json:"weight"`
	Price     decimal.Decimal `json:"price"`
	Package   *PackageType    `json:"package_type,omitempty"`
}

// Order — представление заказа с текущим статусом.
type Order struct {
	OrderCreate
	Status     OrderStatus `json:"status"`
	LastUpdate time.Time   `json:"last_update,omitempty"`
}

// StatusChange — параметр для отслеживания изменения статуса/истории статусов заказа.
type StatusChange struct {
	OrderID OrderID
	Status  OrderStatus
	Time    time.Time
}

// ListOptions — параметры получения списка заказов.
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

// IsExpired — проверяет, истёк ли срок хранения заказа.
func (o *Order) IsExpired(now time.Time) bool {
	return time.Time(o.ExpiresAt).Before(now)
}

// IsExpiredForReturn — проверяет, истёк ли срок возврата заказа.
func (o *Order) IsExpiredForReturn(now time.Time) bool {
	twoDaysAfterExpiry := o.ExpiresAt.Add(48 * time.Hour)
	return now.After(twoDaysAfterExpiry)
}

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusAccepted:
		return "accepted"
	case OrderStatusIssued:
		return "issued"
	case OrderStatusReturned:
		return "returned"
	case OrderStatusDeleted:
		return "deleted"
	default:
		return "unspecified"
	}
}
