package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *Storage) GetOrder(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
	const query = `
		SELECT order_id, user_id, status, expires_at, weight, price, package_type
		FROM orders
		WHERE order_id = $1
	`

	var (
		order   pvz.Order
		pkgType *int32
	)

	err := s.db.QueryRow(ctx, query, orderID).Scan(
		&order.OrderID,
		&order.UserID,
		&order.Status,
		&order.ExpiresAt,
		&order.Weight,
		&order.Price,
		&pkgType,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if pkgType != nil {
		pkg := pvz.PackageType(*pkgType)
		order.Package = &pkg
	}

	return &order, nil
}
