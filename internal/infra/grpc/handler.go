package grpc

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

type OrdersHandler interface {
	Accept(ctx context.Context, order pvz.OrderCreate) error
	Return(ctx context.Context, input app.ReturnOrderInput) error
	Process(ctx context.Context, orders app.ProcessOrdersInput) (app.ProcessOrdersOutput, error)
	ListOrders(ctx context.Context, opts pvz.ListOptions) (*app.OrdersListOutput, error)
	ListReturns(ctx context.Context, opts app.ListReturnsInput) (*app.ReturnsListOutput, error)
	OrderHistory(ctx context.Context) (*app.OrderHistoryList, error)
	StatusHistory(ctx context.Context, orderID pvz.OrderID) ([]*pvz.StatusChange, error)
}
type OrdersServer struct {
	pkg.UnimplementedOrdersServiceServer
	ordersHandler OrdersHandler
}

func NewOrdersServer(handler OrdersHandler) *OrdersServer {
	return &OrdersServer{
		ordersHandler: handler,
	}
}
