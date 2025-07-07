package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
)

func (h *FileHandler) ResizeWorkerPoolCommand(cmd *cobra.Command, _ []string) error {
	n, err := cmd.Flags().GetInt("workers")
	if err != nil || n < 1 {
		return fmt.Errorf("invalid workers count: %v", err)
	}
	ctx := cmd.Context()
	_, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return fmt.Errorf("workerpool not found in context")
	}
	return h.commandsHandler.ResizeWorkerPool(ctx, n)
}
