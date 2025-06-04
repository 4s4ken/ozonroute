package cli

import (
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// OrdersHistoryCommand функция для выполнения команды orders-history
func (h *FileHandler) OrdersHistoryCommand() error {
	history, err := h.commandsHandler.OrderHistory()
	if err != nil {
		return mapError(err)
	}

	for _, ord := range history {
		fmt.Printf(
			"HISTORY: %d %s %s\n", ord.Orderid, ord.Status.String(), pvz.MapTimeToString(ord.Time))
	}

	return nil
}
