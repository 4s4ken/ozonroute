package usecase

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"time"
)

func (s *StorageService) AcceptOrderCommand(cmd *cobra.Command, args []string) error {
	orderID, err := cmd.Flags().GetUint64("order-id")
	userID, err1 := cmd.Flags().GetUint64("user-id")
	expiresStr, _ := cmd.Flags().GetString("expires")

	exp, err2 := time.Parse("2006.01.02", expiresStr)

	if err != nil || err1 != nil || err2 != nil {
		return multierr.Combine(err, err2, err2)
	}

	err = s.Store.Accept(orderID, userID, exp)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Printf("ORDER_ACCEPTED: %d\n", orderID)
	return nil
}
