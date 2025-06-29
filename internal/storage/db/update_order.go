package db

import (
	"context"
	"fmt"

	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *Storage) UpdateOrder(ctx context.Context, order *pvz.Order) (err error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	const deleteOrder = `
		DELETE FROM orders
		WHERE order_id = $1
	`
	commandTag, err := tx.Exec(ctx, deleteOrder, order.OrderID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return storageErrors.EntityNotFoundError("Order", fmt.Sprint(order.OrderID))
	}

	const insertHistory = `
		INSERT INTO order_status_history (order_id, status, changed_at)
		VALUES ($1, $2, NOW())
	`
	_, err = tx.Exec(ctx, insertHistory, order.OrderID, order.Status)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
