package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (s *CommandService) OrdersHistoryCommand(cmd *cobra.Command, args []string) error {
	history, err := s.CommandsHandler.OrderHistory()
	if err != nil {
		return fmt.Errorf("CommandsHandler.OrderHistory: %w", err)
	}

	for _, ord := range history {
		fmt.Printf(
			"HISTORY: %d %s %s\n", ord.Orderid, ord.Status.String(), ord.Time.Format("2006-01-02"))
	}

	return nil
}
