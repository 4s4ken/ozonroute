package cli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"strconv"
	"strings"
)

// ProcessOrdersCommand функция для выполнения команды process-orders
func (h *FileHandler) ProcessOrdersCommand(cmd *cobra.Command, _ []string) error {
	ctx := context.Background()
	userID, err := cmd.Flags().GetUint64("user-id")
	if err != nil {
		return validationError(fmt.Errorf("invalid value of --user-id: %w", err))
	}
	if userID == 0 {
		return validationError(fmt.Errorf("value of --user-id cant equeal 0"))
	}

	action, err := cmd.Flags().GetString("action")
	if err != nil {
		return validationError(fmt.Errorf("invalid value of --action: %w", err))
	}
	orderIDsStr, err := cmd.Flags().GetString("order-ids")
	if err != nil {
		return validationError(fmt.Errorf("invalid value of --order-ids: %w", err))
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
