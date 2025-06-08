package pvz

import (
	"time"
)

type (
	PackageType int
	OrderStatus int
	ActionType  int
	UserID      uint64
	OrderID     uint64
	OrderIDS    []OrderID
)

const (
	PackageTypeUnspecified PackageType = iota
	PackageTypeBag
	PackageTypeBox
	PackageTypeTape
	PackageTypeBagTape
	PackageTypeBoxTape
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

type AcceptOrderInput struct {
	OrderID   uint64
	UserID    uint64
	ExpiresAt time.Time
	Package   *PackageType
	Weight    float32
	Price     float32
}

type ReturnOrderInput struct {
	OrderID OrderID
}

type ProcessOrdersInput struct {
	UserID   UserID
	Action   ActionType
	OrderIDs []OrderID
}

type ListOrdersInput struct {
	UserID     UserID
	InPVZ      bool
	LastN      *uint64
	Pagination *PaginationInput
}

type PaginationInput struct {
	Page        uint64
	CountOnPage uint64
}

type ListReturnsInput struct {
	Pagination *PaginationInput
}

type GetHistoryInput struct {
	Pagination *PaginationInput
}

type Order struct {
	OrderID    OrderID
	UserID     UserID
	Status     OrderStatus
	ExpiresAt  time.Time
	Weight     float32
	TotalPrice float32
	Package    PackageType
}
type ProcessResult struct {
	Processed []uint64
	Errors    []uint64
}

type OrdersList struct {
	Orders []Order
	Total  int32
}

type ReturnsList struct {
	Returns []Order
}

type OrderHistory struct {
	OrderID   uint64
	Status    OrderStatus
	CreatedAt time.Time
}

type OrderHistoryList struct {
	History []OrderHistory
}

type ImportOrdersInput struct {
	Orders []AcceptOrderInput
}

type ImportResult struct {
	Imported int32
	Errors   []uint64
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
