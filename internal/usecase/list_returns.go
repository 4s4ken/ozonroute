package usecase

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *StorageService) ListReturnsCommand(cmd *cobra.Command, args []string) error {
	pagePtr := (*uint64)(nil)
	if cmd.Flags().Changed("page") {
		page, _ := cmd.Flags().GetUint64("page")
		pagePtr = &page
	}

	limitPtr := (*uint64)(nil)
	if cmd.Flags().Changed("limit") {
		limit, _ := cmd.Flags().GetUint64("limit")
		limitPtr = &limit
	}

	options := pvz.ListOrdersOptions{
		Page:  pagePtr,
		Limit: limitPtr,
	}

	orders, err := s.Store.ListReturns(options)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, order := range orders {
		fmt.Printf("RETURN: %d %d %s\n", order.Orderid, order.Userid, order.ReturnedAt.Format("2006.01.02"))
	}

	if pagePtr != nil && limitPtr != nil {
		fmt.Printf("PAGE: %d LIMIT: %d\n", *pagePtr, *limitPtr)
	}

	return nil
}
