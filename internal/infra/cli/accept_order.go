package cli

import (
	"fmt"
	"github.com/shopspring/decimal"
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

	pkg := pvz.NoPackaging
	if cmd.Flags().Changed("package") {
		pkgStr, err := cmd.Flags().GetString("package")
		if err != nil {
			return validationError(err)
		}
		switch pkgStr {
		case string(pvz.NoPackaging), string(pvz.Bag), string(pvz.Box), string(pvz.Film), string(pvz.BagAndFilm), string(pvz.BoxAndFilm):
			pkg = pvz.PackageType(pkgStr)
		default:
			return invalidPackageError(fmt.Errorf("неизвестный тип упаковки: %s", pkgStr))
		}
	}

	totalPrice, err := h.commandsHandler.Accept(nil, pvz.OrderID(orderID), pvz.UserID(userID), exp, pvz.OrderWeight(ordWeight), pvz.OrderPrice(ordPrice), pkg)
	if err != nil {
		return mapError(err)
	}

	fmt.Printf("ORDER_ACCEPTED: %d\nPACKAGE: %s\nTOTAL_PRICE: %s\n", orderID, pkg, totalPrice)
	return nil
}
