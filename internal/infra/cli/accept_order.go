package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"time"
)

func (h *FileHandler) AcceptOrderCommand(cmd *cobra.Command, _ []string) error {
	orderID, err := cmd.Flags().GetUint64("order-id")
	if err != nil {
		return err
	}
	userID, err := cmd.Flags().GetUint64("user-id")
	if err != nil {
		return err
	}
	expiresStr, _ := cmd.Flags().GetString("expires")
	exp, err := time.Parse("2006-01-02", expiresStr)
	if err != nil {
		return err
	}

	err = h.commandsHandler.Accept(nil, pvz.OrderID(orderID), pvz.UserID(userID), exp)
	if err != nil {
		return mapError(err)
	}
	fmt.Printf("ORDER_ACCEPTED: %d\n", orderID)
	return nil
}
