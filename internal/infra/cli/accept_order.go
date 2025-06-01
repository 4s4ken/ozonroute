package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"time"
)

func (h *FileHandler) AcceptOrderCommand(cmd *cobra.Command, _ []string) error {
	orderID, err := cmd.Flags().GetUint64("order-id")
	if err != nil || orderID == 0 {
		return validationError(err)
	}
	userID, err := cmd.Flags().GetUint64("user-id")
	if err != nil || userID == 0 {
		return validationError(err)
	}
	expiresStr, _ := cmd.Flags().GetString("expires")
	exp, err := time.Parse(datetimeFormat, expiresStr)
	if err != nil {
		return validationError(err)
	}

	err = h.commandsHandler.Accept(nil, pvz.OrderID(orderID), pvz.UserID(userID), exp)
	if err != nil {
		return mapError(err)
	}
	fmt.Printf("ORDER_ACCEPTED: %d\n", orderID)
	return nil
}
