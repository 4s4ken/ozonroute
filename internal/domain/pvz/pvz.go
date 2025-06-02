package pvz

import (
	"github.com/shopspring/decimal"
	"time"
)

type (
	OrderID            uint64
	UserID             uint64
	OrderStatus        int
	OrderStatusChanges bool
	OrderWeight        float64
	OrderPrice         = decimal.Decimal
	PackageType        string
)

// Возможные значения типа упаковки.
const (
	NoPackaging PackageType = ""
	Bag                     = "bag"
	Box                     = "box"
	Film                    = "film"
	BagAndFilm              = "bag+film"
	BoxAndFilm              = "box+film"
)

// Возможные статусы заказа.
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
		return ""
	}
}

// ListOrdersOptions задаёт параметры для списка заказов.
type ListOrdersOptions struct {
	UserID uint64
	InPvz  bool
	Last   *uint64
	Page   *uint64
	Limit  *uint64
}

// OrderInfo представляет информацию о заказе.
type OrderInfo struct {
	Orderid       OrderID            `json:"order_id,"`
	Userid        UserID             `json:"user_id"`
	Time          time.Time          `json:"date"`
	Status        OrderStatus        `json:"status"`
	ReturnedAt    time.Time          `json:"returned_at,omitempty"`
	StatusChanged OrderStatusChanges `json:"status_changes"`
	Weight        OrderWeight        `json:"weight"`
	Price         decimal.Decimal    `json:"price"`
	Package       *PackageType       `json:"package_type,omitempty"`
}

// IsExpired проверяет, истёк ли срок хранения заказа.
func (o *OrderInfo) IsExpired(now time.Time) bool {
	return time.Time(o.Time).Before(now)
}

// IsExpiredForReturn проверяет, истёк ли срок возврата заказа.
func (o *OrderInfo) IsExpiredForReturn(now time.Time) bool {
	twoDaysAfterExpiry := o.Time.Add(48 * time.Hour)
	return now.After(twoDaysAfterExpiry)
}

func (p PackageType) String() string {
	return string(p)
}
