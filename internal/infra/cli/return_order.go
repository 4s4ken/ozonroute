package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (s *CommandService) ReturnOrderCommand(cmd *cobra.Command, args []string) error {
	orderID, err := cmd.Flags().GetUint64("order-id")
	if err != nil {
		return fmt.Errorf("CommandsHandler.Return: %w", err)
	}

	err = s.CommandsHandler.Return(orderID)
	if err != nil {
		return fmt.Errorf("%w", err)

	}
	fmt.Printf("ORDER_RETURNED: %d\n", orderID)
	return nil
}
