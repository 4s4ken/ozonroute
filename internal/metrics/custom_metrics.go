package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	OrdersIssued = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_issued",
			Help: "Количество выданных заказов",
		},
	)
	OrdersAccepted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_accepted",
			Help: "Количество принятых заказов",
		},
	)
)

func Register() {
	prometheus.MustRegister(OrdersIssued)
	prometheus.MustRegister(OrdersAccepted)
}
