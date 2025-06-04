package cli

import (
	"fmt"
)

// OrdersHistoryCommand функция для выполнения команды orders-history
func (h *FileHandler) OrdersHistoryCommand() error {
	history, err := h.commandsHandler.OrderHistory()
	if err != nil {
		return mapError(err)
	}

	for _, ord := range history {
		fmt.Printf(
			"HISTORY: %d %s %s\n", ord.Orderid, ord.Status.String(), ord.Time.Format(datetimeFormat))
	}

	return nil
}
