package app_test

import (
	"context"
	"errors"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/mocks"
)

func TestService_ListOrders(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var (
		page  uint64 = 2
		limit uint64 = 5
		lastN uint64 = 10
	)

	tests := []struct {
		name        string
		mockSetup   func(repo *mocks.OrderRepositoryMock, ctx context.Context, input pvz.ListOptions)
		input       pvz.ListOptions
		wantOrders  []pvz.Order
		wantTotal   int32
		wantErr     bool
		errContains string
	}{
		{
			name: "success empty list",
			mockSetup: func(repo *mocks.OrderRepositoryMock, ctx context.Context, input pvz.ListOptions) {
				repo.ListOrdersMock.Expect(ctx, input).Return([]*pvz.Order{}, 0, nil)
			},
			input:      pvz.ListOptions{},
			wantOrders: []pvz.Order{},
			wantTotal:  0,
		},
		{
			name: "success pagination list",
			mockSetup: func(repo *mocks.OrderRepositoryMock, ctx context.Context, input pvz.ListOptions) {
				orders := []*pvz.Order{
					{
						OrderCreate: pvz.OrderCreate{OrderID: pvz.OrderID(1), UserID: pvz.UserID(input.UserID)},
						Status:      pvz.OrderStatusAccepted,
					},
					{
						OrderCreate: pvz.OrderCreate{OrderID: pvz.OrderID(2), UserID: pvz.UserID(input.UserID)},
						Status:      pvz.OrderStatusIssued,
					},
				}
				repo.ListOrdersMock.Expect(ctx, input).Return(orders, 2, nil)
			},
			input: pvz.ListOptions{
				UserID: 42,
				InPVZ:  true,
				LastN:  &lastN,
				Pagination: &pvz.Pagination{
					Page:  &page,
					Limit: &limit,
				},
			},
			wantOrders: []pvz.Order{
				{OrderCreate: pvz.OrderCreate{OrderID: pvz.OrderID(1), UserID: 42}, Status: pvz.OrderStatusAccepted},
				{OrderCreate: pvz.OrderCreate{OrderID: pvz.OrderID(2), UserID: 42}, Status: pvz.OrderStatusIssued},
			},
			wantTotal: 2,
		},
		{
			name: "ошибка репозитория",
			mockSetup: func(repo *mocks.OrderRepositoryMock, ctx context.Context, input pvz.ListOptions) {
				repo.ListOrdersMock.Expect(ctx, input).Return(nil, 0, errors.New("some error"))
			},
			input:       pvz.ListOptions{},
			wantErr:     true,
			errContains: "some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			repo := mocks.NewOrderRepositoryMock(mc)
			if tt.mockSetup != nil {
				tt.mockSetup(repo, ctx, tt.input)
			}

			wp := workerpool.NewWorkerPool(context.Background(), 1)
			svc := app.NewService(repo, nil, wp)

			out, err := svc.ListOrders(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, out)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, out)
			assert.Equal(t, tt.wantTotal, out.Total)
			assert.Equal(t, tt.wantOrders, out.Orders)
		})
	}
}
