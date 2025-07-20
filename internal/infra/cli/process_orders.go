package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"strconv"
	"strings"
)

// ProcessOrdersCommand — обрабатывает CLI команду на обработку заказа.
func (h *FileHandler) ProcessOrdersCommand(cmd *cobra.Command, _ []string) error {
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
	var orderIDs []uint64
	for _, idStr := range ids {
		id, err3 := strconv.Atoi(strings.TrimSpace(idStr))
		if err3 != nil {
			return fmt.Errorf("ERROR: %s is not a valid order ID", idStr)
		}
		orderIDs = append(orderIDs, uint64(id))
	}

	input := app.ProcessOrdersInput{
		UserID:   userID,
		Action:   ActionTypeFromString(action),
		OrderIDs: orderIDs,
	}

	ctx := cmd.Context()
	_, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return fmt.Errorf("workerpool not found in context")
	}
	result, err := h.commandsHandler.Process(ctx, input)
	if err != nil {
		return mapError(err)
	}

	// Вывод успешных и ошибочных обработок
	for _, id := range result.Processed {
		fmt.Printf("PROCESSED: %d\n", id)
	}
	for _, id := range result.Errors {
		fmt.Printf("ERROR: %d\n", id)
	}
	return nil
}
