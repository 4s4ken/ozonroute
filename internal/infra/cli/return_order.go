package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// ReturnOrderCommand функция для выполнения команды return-order
func (h *FileHandler) ReturnOrderCommand(cmd *cobra.Command, _ []string) error {
	orderID, err := cmd.Flags().GetUint64("order-id")
	if err != nil {
		return mapError(err)
	}

	input := app.ReturnOrderInput{
		OrderID: pvz.OrderID(orderID),
	}

	err = h.commandsHandler.Return(nil, input)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
