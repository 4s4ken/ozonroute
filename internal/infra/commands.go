package infra

import (
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type StorageService struct {
	Store Storage
}

type Storage interface {
	Accept(orderID uint64, userID uint64, exp time.Time) error
	Return(orderID uint64) error
	Process(userID uint64, act string, orderIDs []pvz.OrderID) error
	ListOrders(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	ListReturns(typeOfList pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error)
	OrderHistory() ([]*pvz.OrderInfo, error)
	ImportOrders(orders string) (int, error)
	ScrollOrders(userID uint64, limit uint, lastID *pvz.OrderID) ([]*pvz.OrderInfo, *pvz.OrderID, error)
}

func InitCommands(st *StorageService) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "pvz",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "help",
			Short: "h",
			RunE:  HelpCommand,
		},
		&cobra.Command{
			Use:   "accept-order",
			Short: "ao",
			RunE:  st.AcceptOrderCommand,
		},
		&cobra.Command{
			Use:   "return-order",
			Short: "ro",
			RunE:  st.ReturnOrderCommand,
		},
		&cobra.Command{
			Use:   "process-orders",
			Short: "po",
			RunE:  st.ProcessOrdersCommand,
		},
		&cobra.Command{
			Use:   "list-orders",
			Short: "lo",
			RunE:  st.ListOrdersCommand,
		},
		&cobra.Command{
			Use:   "list-returns",
			Short: "lr",
			RunE:  st.ListReturnsCommand,
		},
		&cobra.Command{
			Use:   "order-history",
			Short: "oh",
			RunE:  st.OrdersHistoryCommand,
		},
		&cobra.Command{
			Use:   "import-orders",
			Short: "io",
			RunE:  st.ImportOrdersCommand,
		},
		&cobra.Command{
			Use:   "scroll-orders",
			Short: "so",
			RunE:  st.ScrollOrders,
		},
	)

	addFlags(rootCmd)
	return rootCmd
}

func addFlags(root *cobra.Command) {
	for _, cmd := range root.Commands() {
		switch cmd.Use {
		case "accept-order":
			cmd.Flags().Uint64("order-id", 0, "Order ID")
			cmd.Flags().Uint64("user-id", 0, "User ID")
			cmd.Flags().String("expires", "", "Date of expiration")
		case "return-order":
			cmd.Flags().Uint64("order-id", 0, "Order ID")
		case "process-orders":
			cmd.Flags().Uint64("user-id", 0, "User ID")
			cmd.Flags().String("action", "", "Actions: issue, return")
			cmd.Flags().String("order-ids", "", "Orders ID")
		case "list-orders":
			cmd.Flags().Uint64("user-id", 0, "User ID")
			cmd.Flags().Bool("in-pvz", false, "In PVZ")
			cmd.Flags().Uint64("last", 0, "Last N orders")
			cmd.Flags().Uint64("page", 0, "Page of list")
			cmd.Flags().Uint64("limit", 0, "Limit of page")
		case "list-returns":
			cmd.Flags().Uint64("page", 0, "Page")
			cmd.Flags().Uint64("limit", 0, "Limit")
		case "import-orders":
			cmd.Flags().String("file", "", "Path to JSON file")
		case "scroll-orders":
			cmd.Flags().Uint64("user-id", 0, "User ID")
		}
	}
}
