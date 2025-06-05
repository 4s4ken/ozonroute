package pvz

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"sort"
)

func (s *OrdersService) GetHistory(ctx context.Context, req *pvz.GetHistoryRequest) (*pvz.OrderHistoryList, error) {
	var history []*pvz.OrderHistory
	for _, order := range s.storage {
		history = append(history, &pvz.OrderHistory{
			OrderId:   order.OrderId,
			Status:    order.Status,
			CreatedAt: order.ExpiresAt,
		})
	}
	sort.Slice(history, func(i, j int) bool {
		return history[i].CreatedAt.AsTime().After(history[j].CreatedAt.AsTime())
	})

	if req.Pagination == nil {
		return &pvz.OrderHistoryList{History: history}, nil
	}

	page := uint32(1)
	count := uint32(10)
	if req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}
	if req.Pagination.CountOnPage > 0 {
		count = req.Pagination.CountOnPage
	}

	start := (page - 1) * count
	if start >= uint32(len(history)) {
		return &pvz.OrderHistoryList{History: []*pvz.OrderHistory{}}, nil
	}

	end := start + count
	if end > uint32(len(history)) {
		end = uint32(len(history))
	}

	return &pvz.OrderHistoryList{
		History: history[start:end],
	}, nil
}
