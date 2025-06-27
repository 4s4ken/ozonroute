package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/mocks"
)

func TestService_Accept(t *testing.T) {
	t.Parallel()

	pkgBox := pvz.PackageTypeBox
	pkgUnknown := pvz.PackageType(99)
	now := time.Date(2025, 6, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		input     pvz.OrderCreate
		mockSetup func(repo *mocks.OrderRepositoryMock, ctx context.Context, input pvz.OrderCreate)
		wantErr   bool
		errSubstr string
	}{
		{
			name: "success",
			input: pvz.OrderCreate{
				OrderID:   123,
				UserID:    123,
				Weight:    10,
				Price:     decimal.NewFromInt(100),
				ExpiresAt: now.Add(24 * time.Hour),
				Package:   &pkgBox,
			},
			mockSetup: func(repo *mocks.OrderRepositoryMock, ctx context.Context, input pvz.OrderCreate) {
				repo.GetOrderMock.Expect(ctx, input.OrderID).Return(nil, nil)
				repo.SaveOrderMock.Set(func(ctx context.Context, order *pvz.Order) error {
					assert.Equal(t, input.OrderID, order.OrderID)
					assert.Equal(t, input.UserID, order.UserID)
					assert.Equal(t, pvz.OrderStatusAccepted, order.Status)
					expectedPrice := input.Price.Add(decimal.NewFromInt(20))
					assert.True(t, order.Price.Equal(expectedPrice), "Price should be equal total price")
					return nil
				})
			},
			wantErr: false,
		},

		{
			name: "expired order",
			input: pvz.OrderCreate{
				OrderID:   122,
				UserID:    122,
				Weight:    5,
				Price:     decimal.NewFromInt(100),
				ExpiresAt: now.Add(-1 * time.Hour),
				Package:   &pkgBox,
			},
			mockSetup: func(repo *mocks.OrderRepositoryMock, ctx context.Context, input pvz.OrderCreate) {

			},
			wantErr:   true,
			errSubstr: "Ошибка валидации",
		},

		{
			name: "order already exists",
			input: pvz.OrderCreate{
				OrderID:   121,
				UserID:    121,
				Weight:    5,
				Price:     decimal.NewFromInt(100),
				ExpiresAt: now.Add(24 * time.Hour),
				Package:   &pkgBox,
			},
			mockSetup: func(repo *mocks.OrderRepositoryMock, ctx context.Context, input pvz.OrderCreate) {
				repo.GetOrderMock.Expect(ctx, input.OrderID).Return(&pvz.Order{}, nil)
			},
			wantErr:   true,
			errSubstr: "уже существует",
		},
		{
			name: "invalid package type",
			input: pvz.OrderCreate{
				OrderID:   120,
				UserID:    120,
				Weight:    5,
				Price:     decimal.NewFromInt(100),
				ExpiresAt: now.Add(24 * time.Hour),
				Package:   &pkgUnknown,
			},
			mockSetup: func(repo *mocks.OrderRepositoryMock, ctx context.Context, input pvz.OrderCreate) {
				repo.GetOrderMock.Set(func(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
					return nil, nil
				})

			},
			wantErr:   true,
			errSubstr: "Ошибка при выборе упаковки",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			ctx := context.Background()
			repo := mocks.NewOrderRepositoryMock(mc)

			if tt.mockSetup != nil {
				tt.mockSetup(repo, ctx, tt.input)
			}

			svc := app.NewService(repo, func() time.Time { return now })
			err := svc.Accept(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
