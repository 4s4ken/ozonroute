package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

// ImportOrdersCommand функция для выполнения команды import-orders
func (h *FileHandler) ImportOrdersCommand(cmd *cobra.Command, _ []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	imported, err := h.commandsHandler.ImportOrders(nil, filePath)
	if err != nil {
		return mapError(err)
	}

	fmt.Printf("IMPORTED: %d\n", imported)
	return nil
}
