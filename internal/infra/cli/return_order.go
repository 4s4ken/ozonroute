package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
)

// ReturnOrderCommand — обрабатывает CLI команду на возвращение заказа.
func (h *FileHandler) ReturnOrderCommand(cmd *cobra.Command, _ []string) error {
	orderID, err := cmd.Flags().GetUint64("order-id")
	if err != nil {
		return mapError(err)
	}

	input := app.ReturnOrderInput{
		OrderID: pvz.OrderID(orderID),
	}
	ctx := cmd.Context()
	_, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return fmt.Errorf("workerpool not found in context")
	}
	err = h.commandsHandler.Return(ctx, input)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
