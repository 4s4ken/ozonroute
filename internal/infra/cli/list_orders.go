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

	lastPtr := (*uint64)(nil)
	if cmd.Flags().Changed("last") {
		last, err1 := cmd.Flags().GetUint64("last")
		if err1 != nil {
			return fmt.Errorf("%w", err1)
		}
		lastPtr = &last
	}
	pagePtr := (*uint64)(nil)

	if cmd.Flags().Changed("page") {
		page, err2 := cmd.Flags().GetUint64("page")
		if err2 != nil {
			return fmt.Errorf("%w", err2)
		}
		pagePtr = &page
	}

	limitPtr := (*uint64)(nil)
	if cmd.Flags().Changed("limit") {
		limit, err3 := cmd.Flags().GetUint64("limit")
		if err3 != nil {
			return fmt.Errorf("%w", err3)
		}
		limitPtr = &limit
	}

	options := pvz.ListOrdersOptions{
		UserID: userID,
		InPvz:  inPvz,
		Last:   lastPtr,
		Page:   pagePtr,
		Limit:  limitPtr,
	}

	orders, err := h.commandsHandler.ListOrders(ctx, options)
	if err != nil {
		return mapError(err)
	}

	for _, order := range orders {
		fmt.Printf("ORDER: %d %d %s %s %s %f %s\n",
			order.Orderid,
			order.Userid,
			order.Status.String(),
			order.Time.Format(datetimeFormat),
			order.Package,
			order.Weight,
			order.Price.String(),
		)
	}
	fmt.Printf("TOTAL: %d\n", len(orders))

	return nil
}
