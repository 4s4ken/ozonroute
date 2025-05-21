package usecase

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/domain/pvz"
	"go.uber.org/multierr"
	"strconv"
)

func (s *StorageService) ProcessOrdersCommand(cmd *cobra.Command, args []string) error {
	userID, err := cmd.Flags().GetUint64("user-id")
	action, err1 := cmd.Flags().GetString("action")
	orderIDsStr, err2 := cmd.Flags().GetStringSlice("order-ids")

	if err != nil || err1 != nil || err2 != nil {
		return multierr.Combine(err, err1, err2)
	}
	var orderIDs []pvz.OrderID
	for _, idStr := range orderIDsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("ERROR: %w", idStr)
		}
		orderIDs = append(orderIDs, pvz.OrderID(id))
		fmt.Printf("PROCESSED: %d\n", id)
	}
	err = s.Store.Process(userID, action, orderIDs)
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}
	return nil
}
