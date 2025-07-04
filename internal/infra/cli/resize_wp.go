package cli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

func (h *FileHandler) ResizeWorkerPoolCommand(cmd *cobra.Command, args []string) error {
	n, err := cmd.Flags().GetInt("workers")
	if err != nil || n < 1 {
		return fmt.Errorf("invalid workers count: %v", err)
	}
	return h.commandsHandler.ResizeWorkerPool(context.Background(), n)
}
