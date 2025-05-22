package usecase

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/domain/pvz"
	"os"
	"strings"
)

func (s *StorageService) ScrollOrders(cmd *cobra.Command, args []string) error {
	userID, err := cmd.Flags().GetUint64("user-id")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	limit, err := cmd.Flags().GetUint("limit")
	if err != nil {
		limit = 2
	}

	var lastID *pvz.OrderID
	reader := bufio.NewReader(os.Stdin)

	for {
		orders, nextLastID, err2 := s.Store.ScrollOrders(userID, limit, lastID)
		if err2 != nil {

			return err2
		}
		if len(orders) == 0 {
			fmt.Println("No more orders.")
			break
		}

		for _, order := range orders {
			fmt.Printf("ORDER: %d %d %s %s\n", order.Orderid, order.Userid, order.Status.String(), order.Time.Format("2006.01.02"))
		}

		if nextLastID != nil {
			fmt.Printf("NEXT: %d\n", *nextLastID)
		} else {
			fmt.Println("End of orders. Quitting from scrolling orders.")
			break
		}
		fmt.Println("Type 'next' to continue scroll or 'exit' to quit.")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		switch line {
		case "next":
			lastID = nextLastID
		case "exit":
			return nil
		default:
			fmt.Println("Unknown command. Use 'next' or 'exit'.")
		}
	}

	return nil
}
