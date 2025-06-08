package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// ListReturnsCommand функция для выполнения команды list-returns
func (h *FileHandler) ListReturnsCommand(cmd *cobra.Command, _ []string) error {
	pagePtr := (*uint64)(nil)
	if cmd.Flags().Changed("page") {
		page, err := cmd.Flags().GetUint64("page")
		if err != nil {
			return err
		}
		pagePtr = &page
	}

	limitPtr := (*uint64)(nil)
	if cmd.Flags().Changed("limit") {
		limit, err := cmd.Flags().GetUint64("limit")
		if err != nil {
			return err
		}
		limitPtr = &limit
	}

	options := pvz.ListReturnsInput{}

	if pagePtr != nil && limitPtr != nil {
		options.Pagination = &pvz.PaginationInput{
			Page:        *pagePtr,
			CountOnPage: *limitPtr,
		}
		fmt.Printf("PAGE: %d LIMIT: %d\n", *pagePtr, *limitPtr)
	}

	returnsList, err := h.commandsHandler.ListReturns(nil, options)
	if err != nil {
		return mapError(err)
	}

	for _, order := range returnsList.Returns {
		fmt.Printf("RETURN: %d %d %s\n", order.OrderID, order.UserID, order.ExpiresAt.Format(DatetimeFormat))
	}

	if pagePtr == nil || limitPtr == nil {
		return errors.EntityValidationError("ListReturns", "missing page and/or limit")
	}

	return nil
}
