package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/pflag"
	"gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/domain/pvz"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	//"gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/storage/map_store"
	"gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/storage/json_store"
)

var rootCommand cobra.Command = cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {},
}

type StorageService struct {
	store Storage
}

func HelpCommand(cmd *cobra.Command, args []string) error {
	fmt.Println(`
	accept-order    Принять заказ от курьера
	return-order    Вернуть заказ
	process-orders  Выдать или принять возврат
	list-orders     Получить список заказов
	list-returns    Получить список возвратов
	order-history   Получить историю заказов
	import-orders   Импорт заказов из файла
	`)
	return nil
}

func (s *StorageService) AcceptOrderCommand(cmd *cobra.Command, args []string) error {
	orderID, err := cmd.Flags().GetUint64("order-id")
	userID, err1 := cmd.Flags().GetUint64("user-id")
	expiresStr, _ := cmd.Flags().GetString("expires")

	exp, err2 := time.Parse("2006.01.02", expiresStr)

	if err != nil || err1 != nil || err2 != nil {
		return fmt.Errorf("ERROR: %w", err)
	}

	err = s.store.Accept(orderID, userID, exp)
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}
	fmt.Printf("ORDER_ACCEPTED: %d", orderID)
	return nil
}

func (s *StorageService) ReturnOrderCommand(cmd *cobra.Command, args []string) error {
	orderID, _ := cmd.Flags().GetUint64("order-id")
	err := s.store.Return(orderID)
	if err != nil {
		return fmt.Errorf("ERROR: %w", orderID, err)
	}
	fmt.Printf("ORDER_RETURNED: %d\n", orderID)
	return nil
}

func (s *StorageService) ProcessOrdersCommand(cmd *cobra.Command, args []string) error {
	userID, _ := cmd.Flags().GetUint64("user-id")
	action, _ := cmd.Flags().GetString("action")
	orderIDsStr, _ := cmd.Flags().GetStringSlice("order-ids")

	var orderIDs []pvz.OrderID
	for _, idStr := range orderIDsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("ERROR: %w", idStr)
		}
		orderIDs = append(orderIDs, pvz.OrderID(id))
		fmt.Printf("PROCESSED: %d\n", id)
	}
	err := s.store.Process(userID, action, orderIDs)
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}
	return nil
}

func (s *StorageService) ListOrdersCommand(cmd *cobra.Command, args []string) error {
	userID, _ := cmd.Flags().GetUint64("user-id")

	inPvz := false
	if cmd.Flags().Changed("in-pvz") {
		inPvz, _ = cmd.Flags().GetBool("in-pvz")
	}

	lastPtr := (*uint64)(nil)
	if cmd.Flags().Changed("last") {
		last, _ := cmd.Flags().GetUint64("last")
		lastPtr = &last
	}

	pagePtr := (*uint64)(nil)
	if cmd.Flags().Changed("page") {
		page, _ := cmd.Flags().GetUint64("page")
		pagePtr = &page
	}

	limitPtr := (*uint64)(nil)
	if cmd.Flags().Changed("limit") {
		limit, _ := cmd.Flags().GetUint64("limit")
		limitPtr = &limit
	}

	options := pvz.ListOrdersOptions{
		UserID: userID,
		InPvz:  inPvz,
		Last:   lastPtr,
		Page:   pagePtr,
		Limit:  limitPtr,
	}

	orders, err := s.store.ListOrders(options)
	if err != nil {
		return err
	}

	for _, order := range orders {
		fmt.Printf("ORDER: %d  %s %s\n",
			order.ORDERID, order.Expires.Format("2006.01.02"), order.Status.String())
	}
	fmt.Printf("TOTAL: %d\n", len(orders))
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		err := f.Value.Set(f.DefValue)
		if err != nil {
			return
		}
		f.Changed = false
	})

	return nil
}

func (s *StorageService) ListReturnsCommand(cmd *cobra.Command, args []string) error {
	var (
		pagePtr  *uint64
		limitPtr *uint64
	)

	if cmd.Flags().Changed("page") {
		page, _ := cmd.Flags().GetUint64("page")
		pagePtr = &page
	}
	if cmd.Flags().Changed("limit") {
		limit, _ := cmd.Flags().GetUint64("limit")
		limitPtr = &limit
	}

	options := pvz.ListOrdersOptions{
		Page:  pagePtr,
		Limit: limitPtr,
	}

	orders, err := s.store.ListReturns(options)
	if err != nil {
		return fmt.Errorf("Cannot show list of return: %w", err)
	}

	for _, order := range orders {
		fmt.Printf("RETURN: %d %d %s\n",
			order.ORDERID,
			order.USERID,
			order.ReturnedAt.Format("2006.01.02"))
	}

	if pagePtr != nil && limitPtr != nil {
		fmt.Printf("PAGE: %d LIMIT: %d\n", *pagePtr, *limitPtr)
	}

	return nil
}

func (s *StorageService) OrdersHistoryCommand(cmd *cobra.Command, args []string) error {
	history, err := s.store.OrderHistory()
	if err != nil {
		return fmt.Errorf("Cannot show order history: %w", err)
	}

	for _, ord := range history {
		fmt.Printf(
			"HISTORY: %d %s %s\n",
			ord.OrderID,
			ord.Status.String(),
			ord.Data.Format("2006.01.02"),
		)
	}

	return nil
}

func (s *StorageService) ImportOrdersCommand(cmd *cobra.Command, args []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return fmt.Errorf("Input Error: %w", err)
	}

	imported, err := s.store.ImportOrdersFromFile(filePath)
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}

	fmt.Printf("IMPORTED: %d\n", imported)
	return nil
}

type Storage interface {
	Accept(orderID uint64, userID uint64, exp time.Time) error
	Return(orderID uint64) error
	Process(userID uint64, act string, orderIDs []pvz.OrderID) error
	ListOrders(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	ListReturns(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	OrderHistory() ([]json_store.OrderHistoryEvent, error)
	ImportOrdersFromFile(orders string) (int, error)
}

func main() {
	js, err := json_store.NewJsonStorage("db.json")
	if err != nil {
		log.Fatalf("cannot init json storage: %s", err)
	}

	st := StorageService{
		store: js,
	}

	rootCommand.AddCommand(
		&cobra.Command{
			Use:     "help",
			Short:   "h",
			Example: "",
			Args:    cobra.NoArgs,
			RunE:    HelpCommand,
		},

		&cobra.Command{
			Use:     "accept-order",
			Short:   "ao",
			Example: "",
			RunE:    st.AcceptOrderCommand,
		},

		&cobra.Command{
			Use:     "return-order",
			Short:   "ao",
			Example: "",
			RunE:    st.ReturnOrderCommand,
		},

		&cobra.Command{
			Use:     "process-orders",
			Short:   "po",
			Example: "",
			RunE:    st.ProcessOrdersCommand,
		},

		&cobra.Command{
			Use:     "list-orders",
			Short:   "lo",
			Example: "",
			RunE:    st.ListOrdersCommand,
		},

		&cobra.Command{
			Use:     "list-returns",
			Short:   "lr",
			Example: "",
			RunE:    st.ListReturnsCommand,
		},

		&cobra.Command{
			Use:     "order-history",
			Short:   "oh",
			Example: "",
			RunE:    st.OrdersHistoryCommand,
		},
		&cobra.Command{
			Use:     "import-orders",
			Short:   "oh",
			Example: "",
			RunE:    st.ImportOrdersCommand,
		},
	)
	cmds := rootCommand.Commands()
	for _, cmd := range cmds {
		switch cmd.Use {
		case "accept-order":
			cmd.Flags().Uint64("order-id", 0, "Order ID")
			cmd.Flags().Uint64("user-id", 0, "User ID")
			cmd.Flags().String("expires", "", "Date of expiration")
		case "return-order":
			cmd.Flags().Uint64("order-id", 0, "Order ID")
		case "process-orders":
			cmd.Flags().Uint64("user-id", 0, "User ID")
			cmd.Flags().String("action", "", "Действие: issue, return")
			cmd.Flags().StringSlice("order-ids", []string{}, "ID заказов через запятую")
		case "list-orders":
			cmd.Flags().Uint64("user-id", 0, "User ID")
			cmd.Flags().Bool("in-pvz", false, "Фильтр по ПВЗ")
			cmd.Flags().Uint64("last", 0, "Последние N заказов")
			cmd.Flags().Uint64("page", 0, "Номер страницы")
			cmd.Flags().Uint64("limit", 0, "Размер страницы")
		case "list-returns":
			cmd.Flags().Uint64("page", 0, "Последние N заказов")
			cmd.Flags().Uint64("limit", 0, "Номер страницы")
		case "import-orders":
			cmd.Flags().String("file", "", "path to JSON file with orders")
		}
	}
	fmt.Println("Welcome to the order pick-up point simulation. To view all available commands, type 'help'.")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pvz>")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		rootCommand.SetArgs(parts)

		if err := rootCommand.Execute(); err != nil {
			fmt.Println(os.Stderr, "Command failed, Input Error: %s", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
