package db

import (
	"context"
	"database/sql"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// GetOrder возвращает заказ по orderID.
func (s *Storage) GetOrder(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
	const query = `
		SELECT order_id, user_id, status, expires_at, weight, price, package_type
		FROM orders
		WHERE order_id = $1
	`

	var (
		order   pvz.Order
		pkgType sql.NullInt32
	)
	err := s.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.OrderID,
		&order.UserID,
		&order.Status,
		&order.ExpiresAt,
		&order.Weight,
		&order.Price,
		&pkgType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storageErrors.EntityNotFoundError("Order", fmt.Sprint(orderID))
		}
		return nil, err
	}

	if pkgType.Valid {
		pkg := pvz.PackageType(pkgType.Int32)
		order.Package = &pkg
	}

	return &order, nil
}
