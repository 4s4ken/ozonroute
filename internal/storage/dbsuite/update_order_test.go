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

type UpdateOrderSuite struct {
	Suite
}

func (s *UpdateOrderSuite) TestUpdateOrder_Success() {
	ctx := context.Background()

	order := &pvz.Order{
		OrderCreate: pvz.OrderCreate{
			OrderID:   20,
			UserID:    20,
			Weight:    20,
			Price:     decimal.NewFromInt(100),
			ExpiresAt: time.Now().Add(48 * time.Hour),
			Package:   nil,
		},
		Status: pvz.OrderStatusAccepted,
	}

	err := s.Repo.SaveOrder(ctx, order)
	s.Require().NoError(err)

	order.Status = pvz.OrderStatusIssued
	order.Weight = 25
	order.Price = decimal.NewFromInt(2500)
	order.ExpiresAt = time.Now().Add(72 * time.Hour)

	err = s.Repo.UpdateOrder(ctx, order)
	s.Require().NoError(err)
	
	got, err := s.Repo.GetOrder(ctx, order.OrderID)
	s.Require().NoError(err)
	s.Nil(got)

}

func (s *UpdateOrderSuite) TestUpdateOrder_NotFound() {
	ctx := context.Background()

	order := &pvz.Order{
		OrderCreate: pvz.OrderCreate{
			OrderID:   0,
			UserID:    0,
			Weight:    0,
			Price:     decimal.NewFromInt(100),
			ExpiresAt: time.Now(),
			Package:   nil,
		},
		Status: pvz.OrderStatusAccepted,
	}

	err := s.Repo.UpdateOrder(ctx, order)
	s.Require().Error(err)

	var er *domainErrors.Error
	s.Truef(
		s.ErrorAs(err, &er) && er.Code == domainErrors.ErrorCodeNotFound,
		"ожидалась ошибка EntityNotFound с кодом ErrorCodeNotFound, получено: %v", err,
	)
}

func TestUpdateOrderSuite(t *testing.T) {
	suite.Run(t, new(UpdateOrderSuite))
}
