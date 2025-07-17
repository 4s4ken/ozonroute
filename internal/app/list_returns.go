package app

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"go.opentelemetry.io/otel"
)

type ListReturnsInput struct {
	Pagination *pvz.Pagination
}

type ReturnsListOutput struct {
	Returns []pvz.Order
}

// ListReturns — бизнес-логика команды возврата заказов с пагинацией
func (s *Service) ListReturns(ctx context.Context, input ListReturnsInput) (*ReturnsListOutput, error) {

	tr := otel.Tracer("app.Service")
	ctx, span := tr.Start(ctx, "List_returns")
	defer span.End()

	returned, err := s.repo.ListReturns(ctx, input.Pagination)
	if err != nil {
		return nil, err
	}
	returns := make([]pvz.Order, 0, len(returned))
	for _, o := range returned {
		returns = append(returns, *o)
	}
	return &ReturnsListOutput{
		Returns: returns,
	}, nil
}
