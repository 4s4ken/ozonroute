package cli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"strconv"
	"strings"
)

func (h *FileHandler) ProcessOrdersCommand(cmd *cobra.Command, _ []string) error {
	ctx := context.Background()
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

	err = h.commandsHandler.Process(ctx, pvz.UserID(userID), action, orderIDs)
	if err != nil {
		return mapError(err)
	}

	for _, id := range orderIDs {
		fmt.Printf("PROCESSED: %d\n", id)
	}
	return nil
}
