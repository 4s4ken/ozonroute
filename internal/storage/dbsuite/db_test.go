package dbsuite

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type OrderSuite struct {
	Suite
}

func (s *OrderSuite) TestSaveAndGetOrder() {
	ctx := context.Background()

	order := &pvz.Order{
		OrderCreate: pvz.OrderCreate{
			OrderID:   123,
			UserID:    123,
			Weight:    10,
			Price:     decimal.NewFromInt(100),
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Package:   nil,
		},
		Status: pvz.OrderStatusAccepted,
	}

	err := s.Repo.SaveOrder(ctx, order)
	s.Require().NoError(err)

	got, err := s.Repo.GetOrder(ctx, order.OrderID)
	s.Require().NoError(err)

	s.Equal(order.OrderID, got.OrderID)
	s.Equal(order.UserID, got.UserID)
	s.Equal(order.Status, got.Status)
	s.True(order.Price.Equal(got.Price))

}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}
