package usecase

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *StorageService) ListOrdersCommand(cmd *cobra.Command, args []string) error {
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

	orders, err := s.Store.ListOrders(options)
	if err != nil {
		return err
	}

	for _, order := range orders {
		fmt.Printf("ORDER: %d  %s %s\n", order.Orderid, order.Time.Format("2006.01.02"), order.Status.String())
	}
	fmt.Printf("TOTAL: %d\n", len(orders))
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		err2 := f.Value.Set(f.DefValue)
		if err2 != nil {
			return
		}
		f.Changed = false
	})

	return nil
}
