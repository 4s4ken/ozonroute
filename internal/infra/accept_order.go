package infra

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func (s *StorageService) AcceptOrderCommand(cmd *cobra.Command, args []string) error {
	orderID, err := cmd.Flags().GetUint64("order-id")
	if err != nil {
		return err
	}

	userID, err := cmd.Flags().GetUint64("user-id")
	if err != nil {
		return err
	}
	expiresStr, _ := cmd.Flags().GetString("expires")
	exp, err := time.Parse("2006-01-02", expiresStr)
	if err != nil {
		return err
	}
	err = s.Store.Accept(orderID, userID, exp)
	if err != nil {
		return fmt.Errorf("Store.Accept: %w", err)
	}
	fmt.Printf("ORDER_ACCEPTED: %d\n", orderID)
	return nil
}
