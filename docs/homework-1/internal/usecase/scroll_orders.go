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
		orders, nextLastID, err := s.Store.ScrollOrders(userID, limit, lastID)
		if err != nil {

			return err
		}
		if len(orders) == 0 {
			fmt.Println("No more orders.")
			break
		}

		for _, order := range orders {
			fmt.Printf("ORDER: %d %d %s %s\n", order.ORDERID, order.USERID, order.Status.String(), order.Expires.Format("2006.01.02"))
		}

		if nextLastID != nil {
			fmt.Printf("NEXT: %d\n", *nextLastID)
		} else {
			fmt.Println("End of orders.")
			break
		}

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		fmt.Println("Type 'next' to continue scroll or 'exit' to quit")
		if line == "exit" {
			break
		}
		lastID = nextLastID
	}

	return nil
}
