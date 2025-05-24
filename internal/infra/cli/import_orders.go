package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (s *CommandService) ImportOrdersCommand(cmd *cobra.Command, args []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	imported, err := s.CommandsHandler.ImportOrders(filePath)
	if err != nil {
		return fmt.Errorf("CommandsHandler.ImportOrders: %w", err)
	}

	fmt.Printf("IMPORTED: %d\n", imported)
	return nil
}
