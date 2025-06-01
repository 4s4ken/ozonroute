package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func HelpCommand(_ *cobra.Command, _ []string) error {
	fmt.Println(`
	help
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
