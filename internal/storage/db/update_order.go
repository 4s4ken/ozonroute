package db

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// UpdateOrder обновляет заказ и сохраняет данные в файл.
func (s *Storage) UpdateOrder(ctx context.Context, order *pvz.Order) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const updateOrder = `
		UPDATE orders SET status = $1, expires_at = $2, weight = $3, price = $4, package_type = $5
		WHERE order_id = $6
	`
	res, err := tx.ExecContext(ctx, updateOrder,
		order.Status,
		order.ExpiresAt,
		order.Weight,
		order.Price,
		order.Package,
		order.OrderID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return storageErrors.EntityNotFoundError("Order", fmt.Sprint(order.OrderID))
	}

	const insertHistory = `
		INSERT INTO order_status_history (order_id, status, changed_at)
		VALUES ($1, $2, NOW())
	`
	_, err = tx.ExecContext(ctx, insertHistory, order.OrderID, order.Status)
	if err != nil {
		return err
	}

	return tx.Commit()
}
