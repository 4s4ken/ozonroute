package cli

import (
	"github.com/spf13/cobra"
	"time"
)

func (h *FileHandler) AcceptOrderCommand(cmd *cobra.Command, args []string) error {
	// здесь не обрабатываю ошибки для простоты описания, обрабатывать их конечно же надо
	orderID, err := cmd.Flags().GetUint64("order-id")
	userID, err := cmd.Flags().GetUint64("user-id")
	expiresStr, _ := cmd.Flags().GetString("expires")

	exp, err := time.Parse("2006.01.02", expiresStr)

	err = h.commandsHandler.Accept(orderID, userID, exp) // здесь можешь преобразовать в доменные типы

	return nil
}
