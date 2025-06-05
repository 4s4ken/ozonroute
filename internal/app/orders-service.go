package app

import (
	service "gitlab.ozon.dev/gleb_evseev/homework/internal/service/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

type OrdersServer struct {
	pvz.UnimplementedOrdersServiceServer
	ordersService *service.OrdersService
}

func NewOrdersServer(s *service.OrdersService) *OrdersServer {
	return &OrdersServer{
		ordersService: s,
	}
}
