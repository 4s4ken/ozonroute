package app_test

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	ob "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/mocks"
)

func TestService_Process(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 6, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		action        pvz.ActionType
		orderIDs      []uint64
		mockSetup     func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, orderIDs []uint64)
		wantProcessed []uint64
		wantErrors    []uint64
		wantErr       bool
	}{
		{
			name:     "success issue order",
			action:   pvz.ActionTypeIssue,
			orderIDs: []uint64{1},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, orderIDs []uint64) {
				order := &pvz.Order{
					OrderCreate: pvz.OrderCreate{
						OrderID:   pvz.OrderID(orderIDs[0]),
						ExpiresAt: now.Add(2 * time.Hour),
					},
					Status: pvz.OrderStatusAccepted,
				}
				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					assert.Equal(t, pvz.OrderID(orderIDs[0]), orderID)
					return order, nil
				})
				repo.UpdateOrderMock.Set(func(ctx context.Context, o *pvz.Order) error {
					assert.Equal(t, pvz.OrderStatusIssued, o.Status)
					return nil
				})
				outbox.SaveEventMock.Set(func(ctx context.Context, e ob.Event) error {
					assert.Equal(t, "order_issued", string(e.EventType))
					return nil
				})
			},
			wantProcessed: []uint64{1},
			wantErrors:    nil,
		},
		{
			name:     "success return order",
			action:   pvz.ActionTypeReturn,
			orderIDs: []uint64{2},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, orderIDs []uint64) {
				order := &pvz.Order{
					OrderCreate: pvz.OrderCreate{
						OrderID:   pvz.OrderID(orderIDs[0]),
						ExpiresAt: now.Add(2 * time.Hour),
					},
					Status: pvz.OrderStatusIssued,
				}
				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					assert.Equal(t, pvz.OrderID(orderIDs[0]), orderID)
					return order, nil
				})
				repo.UpdateOrderMock.Set(func(ctx context.Context, o *pvz.Order) error {
					assert.Equal(t, pvz.OrderStatusReturned, o.Status)
					return nil
				})
				outbox.SaveEventMock.Set(func(ctx context.Context, e ob.Event) error {
					assert.Contains(t, []string{"order_returned_by_client"}, string(e.EventType))
					return nil
				})
			},
			wantProcessed: []uint64{2},
			wantErrors:    nil,
		},
		{
			name:     "fail issue expired order",
			action:   pvz.ActionTypeIssue,
			orderIDs: []uint64{3},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, orderIDs []uint64) {
				order := &pvz.Order{
					OrderCreate: pvz.OrderCreate{
						OrderID:   pvz.OrderID(orderIDs[0]),
						ExpiresAt: now.Add(-49 * time.Hour),
					},
					Status: pvz.OrderStatusAccepted,
				}
				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					return order, nil
				})
			},
			wantProcessed: nil,
			wantErrors:    []uint64{3},
		},
		{
			name:     "fail return expired order",
			action:   pvz.ActionTypeReturn,
			orderIDs: []uint64{4},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, orderIDs []uint64) {
				order := &pvz.Order{
					OrderCreate: pvz.OrderCreate{
						OrderID:   pvz.OrderID(orderIDs[0]),
						ExpiresAt: now.Add(-49 * time.Hour),
					},
					Status: pvz.OrderStatusAccepted,
				}
				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					assert.Equal(t, pvz.OrderID(orderIDs[0]), orderID)
					return order, nil
				})
			},
			wantProcessed: nil,
			wantErrors:    []uint64{4},
		},
		{
			name:     "fail get order",
			action:   pvz.ActionTypeIssue,
			orderIDs: []uint64{5},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, orderIDs []uint64) {
				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					assert.Equal(t, pvz.OrderID(orderIDs[0]), orderID)
					return nil, assert.AnError
				})
			},
			wantProcessed: nil,
			wantErrors:    []uint64{5},
		},
		{
			name:     "fail update repo",
			action:   pvz.ActionTypeReturn,
			orderIDs: []uint64{6},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, orderIDs []uint64) {
				order := &pvz.Order{
					OrderCreate: pvz.OrderCreate{
						OrderID:   pvz.OrderID(orderIDs[0]),
						ExpiresAt: now.Add(1 * time.Hour),
					},
					Status: pvz.OrderStatusAccepted,
				}

				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					return order, nil

				})
				repo.UpdateOrderMock.Set(func(ctx context.Context, o *pvz.Order) error {
					return assert.AnError
				})
				outbox.SaveEventMock.Set(func(ctx context.Context, e ob.Event) error {
					return nil
				})
			},
			wantProcessed: nil,
			wantErrors:    []uint64{6},
		},
		{
			name:     "invalid action",
			action:   12,
			orderIDs: []uint64{7},
			mockSetup: func(repo *mocks.OrderRepositoryMock, outbox *mocks.OutboxRepositoryMock, ctx context.Context, orderIDs []uint64) {
			},
			wantErr: true,
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

			if tt.mockSetup != nil {
				tt.mockSetup(repo, outbox, ctx, tt.orderIDs)
			}

			svc := app.NewService(repo, outbox, func() time.Time { return now }, wp)

			out, err := svc.Process(ctx, app.ProcessOrdersInput{
				UserID:   77,
				Action:   tt.action,
				OrderIDs: tt.orderIDs,
			})
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.wantProcessed, out.Processed)
			assert.ElementsMatch(t, tt.wantErrors, out.Errors)
		})
	}
}
