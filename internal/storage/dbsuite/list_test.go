package dbsuite

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type ListSuite struct {
	Suite
}

func (s *ListSuite) SetupTest() {
	// Можно очистить таблицы, чтобы не было мешанины с данными
	_, err := s.pool.Exec(context.Background(), "DELETE FROM orders")
	s.Require().NoError(err)
	_, err = s.pool.Exec(context.Background(), "DELETE FROM order_status_history")
	s.Require().NoError(err)
}

func (s *ListSuite) TestListOrders() {
	ctx := context.Background()
	orders := []*pvz.Order{
		{
			OrderCreate: pvz.OrderCreate{
				OrderID:   123,
				UserID:    123,
				Weight:    10,
				Price:     decimal.NewFromInt(100),
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			Status: pvz.OrderStatusAccepted,
		},
		{
			OrderCreate: pvz.OrderCreate{
				OrderID:   12,
				UserID:    12,
				Weight:    10,
				Price:     decimal.NewFromInt(100),
				ExpiresAt: time.Now().Add(48 * time.Hour),
			},
			Status: pvz.OrderStatusReturned,
		},
		{
			OrderCreate: pvz.OrderCreate{
				OrderID:   13,
				UserID:    13,
				Weight:    10,
				Price:     decimal.NewFromInt(100),
				ExpiresAt: time.Now().Add(72 * time.Hour),
			},
			Status: pvz.OrderStatusAccepted,
		},
	}

	for _, o := range orders {
		err := s.Repo.SaveOrder(ctx, o)
		s.Require().NoError(err)
	}

	opts := pvz.ListOptions{
		UserID: 123,
		InPVZ:  true,
	}

	listedOrders, total, err := s.Repo.ListOrders(ctx, opts)
	s.Require().NoError(err)
	s.Equal(int32(1), total)
	s.Len(listedOrders, 1)
	s.Equal(pvz.OrderID(123), listedOrders[0].OrderID)
}

func (s *ListSuite) TestListReturns() {
	ctx := context.Background()

	order := &pvz.Order{
		OrderCreate: pvz.OrderCreate{
			OrderID:   123,
			UserID:    123,
			Weight:    10,
			Price:     decimal.NewFromInt(100),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
		Status: pvz.OrderStatusReturned,
	}

	err := s.Repo.SaveOrder(ctx, order)
	s.Require().NoError(err)

	returns, err := s.Repo.ListReturns(ctx, nil)
	s.Require().NoError(err)
	s.NotEmpty(returns)

	found := false
	for _, r := range returns {
		if r.OrderID == order.OrderID {
			found = true
		}
	}
	s.True(found, "Should be return status 'Returned'")
}

func (s *ListSuite) TestGetHistoryAndStatusHistory() {
	ctx := context.Background()

	order := &pvz.Order{
		OrderCreate: pvz.OrderCreate{
			OrderID:   30,
			UserID:    4,
			Weight:    15,
			Price:     decimal.NewFromInt(1500),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
		Status: pvz.OrderStatusAccepted,
	}

	err := s.Repo.SaveOrder(ctx, order)
	s.Require().NoError(err)
	order.Status = pvz.OrderStatusIssued
	err = s.Repo.UpdateOrder(ctx, order)
	s.Require().NoError(err)

	history, err := s.Repo.GetHistory(ctx)
	s.Require().NoError(err)
	s.NotEmpty(history)

	statusHistory, err := s.Repo.GetStatusHistory(ctx, order.OrderID)
	s.Require().NoError(err)
	s.NotEmpty(statusHistory)

	s.Equal(pvz.OrderStatusIssued, statusHistory[0].Status)
}

func TestListSuite(t *testing.T) {
	suite.Run(t, new(ListSuite))
}
