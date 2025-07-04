package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	berrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"os"
	"strings"
)

// commandsHandler определяет интерфейс бизнес-логики, которую вызывает CLI
type commandsHandler interface {
	Accept(ctx context.Context, order pvz.OrderCreate) error
	Return(ctx context.Context, input app.ReturnOrderInput) error
	Process(ctx context.Context, orders app.ProcessOrdersInput) (app.ProcessOrdersOutput, error)
	ListOrders(ctx context.Context, opts pvz.ListOptions) (*app.OrdersListOutput, error)
	ListReturns(ctx context.Context, opts app.ListReturnsInput) (*app.ReturnsListOutput, error)
	OrderHistory(ctx context.Context) (*app.OrderHistoryList, error)
	StatusHistory(ctx context.Context, orderID pvz.OrderID) ([]*pvz.StatusChange, error)
	ResizeWorkerPool(ctx context.Context, n int) error
}

type FileHandler struct {
	commandsHandler   commandsHandler
	workerPool        *app.WorkerPool
	scanner           *bufio.Scanner
	commandController *cobra.Command
}

func NewFileHandler(commandsHandler commandsHandler, wp *app.WorkerPool) *FileHandler {
	h := &FileHandler{
		commandsHandler: commandsHandler,
		workerPool:      wp,
		scanner:         bufio.NewScanner(os.Stdin),
	}

	h.commandController = InitCommands(h)

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
			if isFlagParsingError(err) {
				_, _ = fmt.Fprintf(os.Stderr, "Command failed, Input Error: %s\n", validationError(err))
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Command failed, Input Error: %s\n", mapError(err))
			}
		}

	}
	if err := h.scanner.Err(); err != nil {
		return fmt.Errorf("scanner.Err: %s", err)
	}
	return nil
}

// isFlagParsingError определяет, является ли ошибка валидации при парсинге флагов.
func isFlagParsingError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "invalid argument") ||
		strings.Contains(msg, "unknown flag") ||
		strings.Contains(msg, "flag needs an argument") ||
		strings.Contains(msg, "expected value")
}

// mapError преобразует бизнес-ошибку в читаемое сообщение об ошибке с кодом
func mapError(err error) error {
	var be *berrors.Error
	if errors.As(err, &be) {
		switch be.Code {
		case berrors.ErrorCodeNotFound:
			return notFoundError(be)
		case berrors.ErrorCodeAlreadyExists:
			return alreadyExistsError(be)
		case berrors.ErrorCodeExpires:
			return expiredError(be)
		case berrors.ErrorCodeValidation:
			return validationError(be)
		case berrors.ErrorCodePackage:
			return invalidPackageError(be)
		case berrors.ErrorCodeWeight:
			return weightTooHeavyError(be)
		default:
			return internalError(be)
		}
	}
	return internalError(err)
}

// Функции-шаблоны, оборачивающие ошибки с префиксом конкретного кода
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
func invalidPackageError(err error) error {
	return fmt.Errorf("INVALID_PACKAGE: %w", err)
}
func weightTooHeavyError(err error) error {
	return fmt.Errorf("WEIGHT_TOO_HEAVY: %w", err)
}
