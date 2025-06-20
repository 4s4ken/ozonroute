package db

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"strings"
	"time"
)

func (s *Storage) ListOrders(ctx context.Context, options pvz.ListOptions) ([]*pvz.Order, int32, error) {
	var (
		whereClauses []string
		args         []interface{}
		argIdx       = 1
	)

	whereClauses = append(whereClauses, fmt.Sprintf("user_id = $%d", argIdx))
	args = append(args, options.UserID)
	argIdx++

	if options.InPVZ {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, pvz.OrderStatusAccepted)
		argIdx++
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	var total int32
	countQuery := "SELECT COUNT(*) FROM orders " + whereSQL
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
	SELECT order_id, user_id, status, expires_at, weight, price, package_type
	FROM orders
	whereSQL
	ORDER BY order_id
`)

	if options.Pagination != nil && options.Pagination.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *options.Pagination.Limit)
	}
	if options.Pagination != nil && options.Pagination.Page != nil && options.Pagination.Limit != nil {
		offset := (*options.Pagination.Page) * (*options.Pagination.Limit)
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*pvz.Order
	for rows.Next() {
		var o pvz.Order
		var pkgType *int32
		if err := rows.Scan(
			&o.OrderID,
			&o.UserID,
			&o.Status,
			&o.ExpiresAt,
			&o.Weight,
			&o.Price,
			&pkgType,
		); err != nil {
			return nil, 0, err
		}
		if pkgType != nil {
			pkg := pvz.PackageType(*pkgType)
			o.Package = &pkg
		}
		orders = append(orders, &o)
	}

	return orders, total, nil
}

func (s *Storage) ListReturns(ctx context.Context, pagination *pvz.Pagination) ([]*pvz.Order, error) {
	query := `
		SELECT order_id, user_id, status, expires_at, weight, price, package_type
		FROM orders
		WHERE status = $1
		ORDER BY expires_at DESC
	`

	args := []interface{}{pvz.OrderStatusReturned}

	if pagination != nil && pagination.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *pagination.Limit)
	}
	if pagination != nil && pagination.Page != nil && pagination.Limit != nil {
		offset := (*pagination.Page) * (*pagination.Limit)
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var returns []*pvz.Order
	for rows.Next() {
		var o pvz.Order
		var pkgType *int32
		if err := rows.Scan(
			&o.OrderID,
			&o.UserID,
			&o.Status,
			&o.ExpiresAt,
			&o.Weight,
			&o.Price,
			&pkgType,
		); err != nil {
			return nil, err
		}
		if pkgType != nil {
			pkg := pvz.PackageType(*pkgType)
			o.Package = &pkg
		}
		returns = append(returns, &o)
	}

	return returns, nil
}

func (s *Storage) GetHistory(ctx context.Context) ([]*pvz.StatusChange, error) {
	query := `
		SELECT order_id, status, changed_at
		FROM order_status_history
		ORDER BY changed_at DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*pvz.StatusChange
	for rows.Next() {
		var h pvz.StatusChange
		if err := rows.Scan(&h.OrderID, &h.Status, &h.Time); err != nil {
			return nil, err
		}
		history = append(history, &h)
	}

	return history, nil
}

func (s *Storage) GetStatusHistory(ctx context.Context, orderID pvz.OrderID) ([]*pvz.StatusChange, error) {
	const query = `
		SELECT status, changed_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY changed_at DESC
	`

	rows, err := s.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*pvz.StatusChange
	for rows.Next() {
		var status int
		var changedAt time.Time
		if err := rows.Scan(&status, &changedAt); err != nil {
			return nil, err
		}
		history = append(history, &pvz.StatusChange{
			OrderID: orderID,
			Status:  pvz.OrderStatus(status),
			Time:    changedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return nil, storageErrors.EntityNotFoundError("Order", fmt.Sprint(orderID))
	}

	return history, nil
}
