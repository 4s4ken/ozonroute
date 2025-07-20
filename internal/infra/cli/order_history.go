package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
)

// OrdersHistoryCommand — обрабатывает CLI команду на вывод списка истории изменений заказов.
func (h *FileHandler) OrdersHistoryCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	_, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return fmt.Errorf("workerpool not found in context")
	}
	history, err := h.commandsHandler.OrderHistory(ctx)
	if err != nil {
		return mapError(err)
	}

	for _, ord := range history.History {
		fmt.Printf(
			"HISTORY: %d %s %s\n", ord.OrderID, MapStatusToString(ord.Status), MapTimeToString(ord.Time))
	}

	return nil
}
