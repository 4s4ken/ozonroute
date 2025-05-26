package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (h *FileHandler) ReturnOrderCommand(cmd *cobra.Command, _ []string) error {
	orderID, err := cmd.Flags().GetUint64("order-id")
	if err != nil {
		return mapError(err)
	}

	err = h.commandsHandler.Return(nil, pvz.OrderID(orderID))
	if err != nil {
		return fmt.Errorf("%w", err)

	}
	fmt.Printf("ORDER_RETURNED: %d\n", orderID)
	return nil
}
