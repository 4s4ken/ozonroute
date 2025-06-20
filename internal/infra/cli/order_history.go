package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

// OrdersHistoryCommand функция для выполнения команды orders-history
func (h *FileHandler) OrdersHistoryCommand(cmd *cobra.Command, args []string) error {
	history, err := h.commandsHandler.OrderHistory(nil)
	if err != nil {
		return mapError(err)
	}

	for _, ord := range history.History {
		fmt.Printf(
			"HISTORY: %d %s %s\n", ord.OrderID, MapStatusToString(ord.Status), MapTimeToString(ord.Time))
	}

	return nil
}
