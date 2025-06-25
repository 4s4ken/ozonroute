package db

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *Storage) SaveOrder(ctx context.Context, order *pvz.Order) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	const insertOrder = `
		INSERT INTO orders (order_id, user_id, status, expires_at, weight, price, package_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.Exec(ctx, insertOrder,
		order.OrderID,
		order.UserID,
		order.Status,
		order.ExpiresAt,
		order.Weight,
		order.Price,
		toDBPackageType(order.Package),
	)
	if err != nil {
		return err
	}

	const insertHistory = `
		INSERT INTO order_status_history (order_id, status, changed_at)
		VALUES ($1, $2, NOW())
	`

	_, err = tx.Exec(ctx, insertHistory, order.OrderID, order.Status)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}
