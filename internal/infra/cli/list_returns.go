package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (h *FileHandler) ListReturnsCommand(cmd *cobra.Command, _ []string) error {
	pagePtr := (*uint64)(nil)
	if cmd.Flags().Changed("page") {
		page, err := cmd.Flags().GetUint64("page")
		if err != nil {
			return err
		}
		pagePtr = &page
	}

	limitPtr := (*uint64)(nil)
	if cmd.Flags().Changed("limit") {
		limit, err := cmd.Flags().GetUint64("limit")
		if err != nil {
			return err
		}
		limitPtr = &limit
	}

	options := pvz.ListOrdersOptions{
		Page:  pagePtr,
		Limit: limitPtr,
	}

	orders, err := h.commandsHandler.ListReturns(nil, options)
	if err != nil {
		return mapError(err)
	}

	for _, order := range orders {
		fmt.Printf("RETURN: %d %d %s\n", order.Orderid, order.Userid, order.ReturnedAt.Format(datetimeFormat))
	}

	if pagePtr != nil && limitPtr != nil {
		fmt.Printf("PAGE: %d LIMIT: %d\n", *pagePtr, *limitPtr)
	}

	return nil
}
