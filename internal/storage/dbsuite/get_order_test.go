package dbsuite

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	domainErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type GetOrderSuite struct {
	Suite
}

func (s *GetOrderSuite) TestGetOrder_Found() {
	ctx := context.Background()

	order := &pvz.Order{
		OrderCreate: pvz.OrderCreate{
			OrderID:   123,
			UserID:    123,
			Weight:    10,
			Price:     decimal.NewFromInt(100),
			ExpiresAt: time.Now().Add(48 * time.Hour),
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
	s.Equal(order.Weight, got.Weight)
	s.WithinDuration(order.ExpiresAt, got.ExpiresAt, time.Second)
	s.Nil(got.Package)
}

func (s *GetOrderSuite) TestGetOrder_NotFound() {
	ctx := context.Background()

	_, err := s.Repo.GetOrder(ctx, 123123)
	s.Require().Error(err)

	var er *domainErrors.Error
	s.Truef(
		s.ErrorAs(err, &er) && er.Code == domainErrors.ErrorCodeNotFound,
		"Waiting EntityNotFound error,got ERROR: %v", err,
	)
}

func TestGetOrderSuite(t *testing.T) {
	suite.Run(t, new(GetOrderSuite))
}
