package cli

import (
	"bufio"
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
	Accept(orderID uint64, userID uint64, exp time.Time) error
	Return(orderID uint64) error
	Process(userID uint64, act string, orderIDs []pvz.OrderID) error
	ListOrders(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	ListReturns(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	OrderHistory() ([]*pvz.OrderInfo, error)
	ImportOrders(orders string) (int, error)
	ScrollOrders(userID uint64, limit uint, lastID *pvz.OrderID) ([]*pvz.OrderInfo, *pvz.OrderID, error)
}

type FileHandler struct {
	commandsHandler   commandsHandler
	scanner           *bufio.Scanner
	commandController *cobra.Command
}

func NewFileHandler(commandsHandler commandsHandler) *FileHandler {
	h := &FileHandler{
		scanner: bufio.NewScanner(os.Stdin),
	}

	h.commandController = InitCommands(h)
	h.commandController.AddCommand(
		&cobra.Command{
			Use:   "help",
			Short: "h",
			RunE:  HelpCommand,
		},
		&cobra.Command{
			Use:   "accept-order",
			Short: "ao",
			RunE:  h.AcceptOrderCommand,
		},
		// и тд
	)

	addFlags(h.commandController)

	return h
}

func (h *FileHandler) Run() error {
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

		if err := h.commandController.Execute(); err != nil { // лучше по всему стеку всегда передавать контекст, в го часто это единственный способ передавать метаданные по коды (например, авторизационные данные, которые могут понадобиться где-нибудь в глубине приложени)
			_, _ = fmt.Fprintf(os.Stderr, "Command failed, Input Error: %s\n", mapError(err)) // здесь маппим бизнесовые ошибки на error-коды, фильтруем что можем показать юзерам, а что нет
		}
	}
	if err := h.scanner.Err(); err != nil {
		return fmt.Errorf("scanner.Err: %s", err)
	}
	return nil
}

func mapError(err error) error {
	if businessError := new(businesserror.Error); errors.As(err, businessError) {
		switch businessError.Code {
		case businesserror.ErrorCodeNotFound:
			// businessError не будет уже содержать весь стек вызовов,
			// она будет содержать только исходную ошибку
			return NotFoundErrorError(businessError)
			// ...
			// и остальные case
		default:
			return InternalError(businessError)
		}
	} else {
		return InternalError(err)
	}
}

func NotFoundErrorError(err error) error {
	return fmt.Errorf("NOT_FOUND_ERROR: %w", err) // здесь можешь использовать свои коды ответов, а не хардкодить строку
	// так же здесь можно реализовать какую нибудь логику которая "спрячет" реализацию и стек вызовов от пользователя
}

func InternalError(err error) error {
	return fmt.Errorf("INTERNAL_ERROR: %w", err) // здесь можешь использовать свои коды ответов, а не хардкодить строку
}
