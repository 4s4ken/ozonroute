package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (h *FileHandler) OrdersHistoryCommand(_ *cobra.Command, _ []string) error {
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
