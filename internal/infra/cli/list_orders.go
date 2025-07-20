package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
)

// ListOrdersCommand — обрабатывает CLI команду на вывод списка заказов.
func (h *FileHandler) ListOrdersCommand(cmd *cobra.Command, _ []string) error {
	userID, err := cmd.Flags().GetUint64("user-id")
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}

	var inPvz bool
	if cmd.Flags().Changed("in-pvz") {
		inPvz, _ = cmd.Flags().GetBool("in-pvz")
	}

	var lastPtr *uint64
	if cmd.Flags().Changed("last") {
		last, err := cmd.Flags().GetUint64("last")
		if err != nil {
			return fmt.Errorf("ERROR: %w", err)
		}
		lastPtr = &last
	}

	var pagePtr, limitPtr *uint64
	if cmd.Flags().Changed("page") {
		page, err := cmd.Flags().GetUint64("page")
		if err != nil {
			return fmt.Errorf("ERROR: %w", err)
		}
		pagePtr = &page
	}
	if cmd.Flags().Changed("limit") {
		limit, err := cmd.Flags().GetUint64("limit")
		if err != nil {
			return fmt.Errorf("ERROR: %w", err)
		}
		limitPtr = &limit
	}

	options := pvz.ListOptions{
		UserID: userID,
		InPVZ:  inPvz,
		LastN:  lastPtr,
		Pagination: &pvz.Pagination{
			Page:  pagePtr,
			Limit: limitPtr,
		},
	}

	ctx := cmd.Context()
	_, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return fmt.Errorf("workerpool not found in context")
	}
	ordersList, err := h.commandsHandler.ListOrders(ctx, options)
	if err != nil {
		return mapError(err)
	}

	for _, order := range ordersList.Orders {
		pkgStr := "none"
		if order.Package != nil {
			pkgStr = string(*order.Package)
		}
		fmt.Printf("ORDER: %d %d %d %s %s %.2f %s\n",
			order.OrderID,
			order.UserID,
			order.Status,
			MapTimeToString(order.ExpiresAt),
			pkgStr,
			order.Weight,
			order.Price.String(),
		)
	}

	fmt.Printf("TOTAL: %d\n", ordersList.Total)
	return nil
}
