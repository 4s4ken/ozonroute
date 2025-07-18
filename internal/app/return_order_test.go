package app_test

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/config"
	ob "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/cache"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/mocks"
)

func TestService_Return(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 6, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		input     app.ReturnOrderInput
		mockSetup func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, input app.ReturnOrderInput)
		wantErr   bool
		errSubstr string
	}{
		{
			name: "success",
			input: app.ReturnOrderInput{
				UserID:  77,
				OrderID: 1,
			},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, input app.ReturnOrderInput) {
				order := &pvz.Order{
					OrderCreate: pvz.OrderCreate{
						OrderID:   input.OrderID,
						UserID:    input.UserID,
						ExpiresAt: now.Add(24 * time.Hour),
					},
					Status: pvz.OrderStatusAccepted,
				}

				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					assert.Equal(t, input.OrderID, orderID)
					return order, nil
				})

				repo.UpdateOrderMock.Set(func(ctx context.Context, o *pvz.Order) error {
					assert.Equal(t, pvz.OrderStatusReturned, o.Status)
					assert.WithinDuration(t, now, o.OrderCreate.ExpiresAt, time.Second)
					return nil
				})
				outbox.SaveEventMock.Set(func(ctx context.Context, e ob.Event) error {
					assert.Equal(t, "order_returned_to_courier", string(e.EventType))
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "order already returned",
			input: app.ReturnOrderInput{
				UserID:  77,
				OrderID: 2,
			},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, input app.ReturnOrderInput) {
				order := &pvz.Order{
					OrderCreate: pvz.OrderCreate{
						OrderID: input.OrderID,
						UserID:  input.UserID,
					},
					Status: pvz.OrderStatusReturned,
				}
				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					return order, nil
				})
			},
			wantErr:   true,
			errSubstr: "Ошибка валидации",
		},
		{
			name: "order expired while accepted",
			input: app.ReturnOrderInput{
				UserID:  77,
				OrderID: 3,
			},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, input app.ReturnOrderInput) {
				order := &pvz.Order{
					OrderCreate: pvz.OrderCreate{
						OrderID:   input.OrderID,
						UserID:    input.UserID,
						ExpiresAt: now.Add(-1 * time.Hour),
					},
					Status: pvz.OrderStatusAccepted,
				}
				repo.GetOrderMock.Set(func(ctx context.Context, OrderID pvz.OrderID) (*pvz.Order, error) {
					return order, nil
				})
			},
			wantErr:   true,
			errSubstr: "просрочена",
		},

		{
			name: "failed get order",
			input: app.ReturnOrderInput{
				UserID:  77,
				OrderID: 4,
			},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, input app.ReturnOrderInput) {
				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					return nil, assert.AnError
				})
			},
			wantErr:   true,
			errSubstr: assert.AnError.Error(),
		},
		{
			name: "failed update repo",
			input: app.ReturnOrderInput{
				UserID:  77,
				OrderID: 5,
			},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, input app.ReturnOrderInput) {
				order := &pvz.Order{
					OrderCreate: pvz.OrderCreate{
						OrderID:   input.OrderID,
						UserID:    input.UserID,
						ExpiresAt: now.Add(1 * time.Hour),
					},
					Status: pvz.OrderStatusAccepted,
				}
				repo.GetOrderMock.Set(func(ctx context.Context, OrderID pvz.OrderID) (*pvz.Order, error) {
					return order, nil
				})
				repo.UpdateOrderMock.Set(func(ctx context.Context, o *pvz.Order) error {
					return assert.AnError
				})
			},
			wantErr:   true,
			errSubstr: "update order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			ctx := context.Background()
			wp := workerpool.NewWorkerPool(ctx, 1)
			ctx = workerpool.WorkerPoolValue(ctx, wp)
			repo := mocks.NewOrderRepositoryMock(mc)
			outbox := mocks.NewOutboxRepositoryMock(mc)
			orderCache := cache.New[string, *pvz.Order](config.CacheCfg)

			if tt.mockSetup != nil {
				tt.mockSetup(repo, outbox, ctx, tt.input)
			}
			svc := app.NewService(repo, outbox, func() time.Time { return now }, wp, orderCache)

			err := svc.Return(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
