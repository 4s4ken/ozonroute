package usecase

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (s *StorageService) ImportOrdersCommand(cmd *cobra.Command, args []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	imported, err := s.Store.ImportOrders(filePath)
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}

	fmt.Printf("IMPORTED: %d\n", imported)
	return nil
}
