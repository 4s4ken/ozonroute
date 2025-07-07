package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
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

	options := app.ListReturnsInput{}

	if pagePtr != nil && limitPtr != nil {
		options = app.ListReturnsInput{
			Pagination: &pvz.Pagination{
				Page:  pagePtr,
				Limit: limitPtr,
			},
		}
		fmt.Printf("PAGE: %d LIMIT: %d\n", *pagePtr, *limitPtr)
	}

	ctx := cmd.Context()
	_, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return fmt.Errorf("workerpool not found in context")
	}

	returnsList, err := h.commandsHandler.ListReturns(ctx, options)
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
