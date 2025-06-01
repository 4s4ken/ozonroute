package cli

import (
	"github.com/spf13/cobra"
)

func InitCommands(st *FileHandler) *cobra.Command {
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
			cmd.Flags().Float64("weight", 0, "Weight of order")
			cmd.Flags().String("price", "0", "Price of order")
			cmd.Flags().String("package", "", "Type of order package")
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
