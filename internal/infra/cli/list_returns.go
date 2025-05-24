package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *CommandService) ListReturnsCommand(cmd *cobra.Command, args []string) error {
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

	orders, err := s.CommandsHandler.ListReturns(options)
	if err != nil {
		return fmt.Errorf("CommandsHandler.LisReturns: %w", err)
	}

	for _, order := range orders {
		fmt.Printf("RETURN: %d %d %s\n", order.Orderid, order.Userid, order.ReturnedAt.Format("2006-01-02"))
	}

	if pagePtr != nil && limitPtr != nil {
		fmt.Printf("PAGE: %d LIMIT: %d\n", *pagePtr, *limitPtr)
	}

	return nil
}
