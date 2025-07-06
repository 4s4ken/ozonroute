package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
)

// OrdersHistoryCommand функция для выполнения команды orders-history
func (h *FileHandler) OrdersHistoryCommand(cmd *cobra.Command, args []string) error {
	ctx := workerpool.WorkerPoolValue(cmd.Context(), h.workerPool)
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
