package cli

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
)

// AcceptOrderCommand — обрабатывает CLI команду на принятия заказа.
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
	exp, err := MapStringToTime(expiresStr)
	if err != nil {
		return validationError(err)
	}

	ordWeight, err := cmd.Flags().GetFloat64("weight")
	if err != nil || ordWeight == 0 {
		return validationError(err)
	}

	priceStr, err := cmd.Flags().GetString("price")
	if err != nil {
		return validationError(err)
	}
	ordPrice, err := decimal.NewFromString(priceStr)
	if err != nil {
		return validationError(err)
	}

	pkgStr, err := cmd.Flags().GetString("package")
	if err != nil {
		return validationError(err)
	}

	var pkg *pvz.PackageType
	if pkgStr != "" {
		pkgVal, err := MapPackageToDomain(pkgStr)
		if err != nil {
			return err
		}
		pkg = &pkgVal
	}
	order := pvz.OrderCreate{
		OrderID:   pvz.OrderID(orderID),
		UserID:    pvz.UserID(userID),
		ExpiresAt: exp,
		Weight:    float32(ordWeight),
		Price:     ordPrice,
		Package:   pkg,
	}

	ctx := cmd.Context()
	_, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return fmt.Errorf("workerpool not found in context")
	}
	err = h.commandsHandler.Accept(ctx, order)
	if err != nil {
		return mapError(err)
	}

	fmt.Printf("ORDER: %d\nPACKAGE: %s\nTOTAL_PRICE: %s\n", orderID, pkgStr, ordPrice)
	return nil
}
