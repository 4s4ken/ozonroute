package cli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// ListOrdersCommand функция для выполнения команды list-orders
func (h *FileHandler) ListOrdersCommand(cmd *cobra.Command, _ []string) error {
	ctx := context.Background()

	userID, err := cmd.Flags().GetUint64("user-id")
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}

	inPvz := false
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

	var pagination *pvz.PaginationInput
	if pagePtr != nil || limitPtr != nil {
		pagination = &pvz.PaginationInput{}
		if pagePtr != nil {
			pagination.Page = *pagePtr
		}
		if limitPtr != nil {
			pagination.CountOnPage = *limitPtr
		}
	}

	options := pvz.ListOrdersInput{
		UserID:     pvz.UserID(userID),
		InPVZ:      inPvz,
		LastN:      lastPtr,
		Pagination: pagination,
	}

	ordersList, err := h.commandsHandler.ListOrders(ctx, options)
	if err != nil {
		return mapError(err)
	}

	for _, order := range ordersList.Orders {
		pkgStr := "none"
		fmt.Printf("ORDER: %d %d %d %s %s %.2f %f\n",
			order.OrderID,
			order.UserID,
			order.Status,
			MapTimeToString(order.ExpiresAt),
			pkgStr,
			order.Weight,
			order.TotalPrice,
		)
	}

	fmt.Printf("TOTAL: %d\n", ordersList.Total)
	return nil
}
