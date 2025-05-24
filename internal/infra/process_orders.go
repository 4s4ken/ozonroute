package infra

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"strconv"
	"strings"
)

func (s *StorageService) ProcessOrdersCommand(cmd *cobra.Command, args []string) error {
	userID, err := cmd.Flags().GetUint64("user-id")
	if err != nil {
		return err
	}
	action, err := cmd.Flags().GetString("action")
	if err != nil {
		return err
	}
	orderIDsStr, err := cmd.Flags().GetString("order-ids")
	if err != nil {
		return err
	}

	ids := strings.Split(orderIDsStr, ",")
	var orderIDs []pvz.OrderID
	for _, idStr := range ids {
		id, err3 := strconv.Atoi(strings.TrimSpace(idStr))
		if err3 != nil {
			return fmt.Errorf("ERROR: %s is not a valid order ID", idStr)
		}
		orderIDs = append(orderIDs, pvz.OrderID(id))
	}

	err = s.Store.Process(userID, action, orderIDs)
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}

	for _, id := range orderIDs {
		fmt.Printf("PROCESSED: %d\n", id)
	}
	return nil
}
