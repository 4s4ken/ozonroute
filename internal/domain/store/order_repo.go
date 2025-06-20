package store

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *pvz.Order) error
	GetOrder(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error)
	UpdateOrder(ctx context.Context, order *pvz.Order) error
	ListOrders(ctx context.Context, options pvz.ListOptions) ([]*pvz.Order, int32, error)
	ListReturns(ctx context.Context, pagination *pvz.Pagination) ([]*pvz.Order, error)
	GetHistory(ctx context.Context) ([]*pvz.StatusChange, error)
	GetStatusHistory(ctx context.Context, orderID pvz.OrderID) ([]*pvz.StatusChange, error)
}
