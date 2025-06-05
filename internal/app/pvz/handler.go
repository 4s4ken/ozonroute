package pvz

import (
	"context"
	desc "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

type Handler struct {
	desc.UnimplementedOrdersServiceServer
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AcceptOrder(ctx context.Context, req *desc.AcceptOrderRequest) (*desc.AcceptOrderResponse, error) {
	err := h.service.AcceptOrder(req.OrderId, req.UserId, req.Ttl)
	if err != nil {
		return nil, mapError(err) // функция, маппящая ошибки бизнес-логики в gRPC
	}
	return &desc.AcceptOrderResponse{Success: true}, nil
}
