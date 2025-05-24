package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"os"
	"strings"
	"time"
)

type commandsHandler interface {
	Accept(ctx context.Context, orderID pvz.OrderID, userID pvz.UserID, exp time.Time) error
	Return(ctx context.Context, orderID pvz.OrderID) error
	Process(ctx context.Context, userID pvz.UserID, act string, orderIDs []pvz.OrderID) error
	ListOrders(ctx context.Context, typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	ListReturns(ctx context.Context, typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	OrderHistory() ([]*pvz.OrderInfo, error)
	ImportOrders(ctx context.Context, orders string) (int, error)
	ScrollOrders(ctx context.Context, userID pvz.UserID, limit uint, lastID *pvz.OrderID) ([]*pvz.OrderInfo, *pvz.OrderID, error)
}

type FileHandler struct {
	commandsHandler   commandsHandler
	scanner           *bufio.Scanner
	commandController *cobra.Command
}

func NewFileHandler(commandsHandler commandsHandler) *FileHandler {
	h := &FileHandler{
		commandsHandler: commandsHandler,
		scanner:         bufio.NewScanner(os.Stdin),
	}

	h.commandController = InitCommands(h)

	//addFlags(h.commandController)

	return h
}

func (h *FileHandler) Run(ctx context.Context) error {

	for {
		fmt.Print(">")
		if !h.scanner.Scan() {
			break
		}
		line := strings.TrimSpace(h.scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		h.commandController.SetArgs(parts)

		if err := h.commandController.Execute(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Command failed, Input Error: %s\n", mapError(err))
		}
	}
	if err := h.scanner.Err(); err != nil {
		return fmt.Errorf("scanner.Err: %s", err)
	}
	return nil
}

func mapError(err error) error {
	var be *businesserror.Error
	if errors.As(err, &be) {
		switch be.Code {
		case businesserror.ErrorCodeNotFound:
			return notFoundError(be)
		case businesserror.ErrorCodeAlreadyExists:
			return alreadyExistsError(be)
		case businesserror.ErrorCodeExpires:
			return expiredError(be)
		case businesserror.ErrorCodeValidation:
			return validationError(be)
		default:
			return internalError(be)
		}
	}
	return internalError(err)
}

func notFoundError(err error) error {
	return fmt.Errorf("NOT_FOUND: %w", err)
}

func alreadyExistsError(err error) error {
	return fmt.Errorf("ALREADY_EXISTS: %w", err)
}

func expiredError(err error) error {
	return fmt.Errorf("EXPIRED: %w", err)
}

func validationError(err error) error {
	return fmt.Errorf("VALIDATION_FAILED: %w", err)
}

func internalError(err error) error {
	return fmt.Errorf("INTERNAL_ERROR: %w", err)
}
