package config

const (
	DBPath       = "./db/ordersDB.json"
	DBConnString = "user=postgres password=postgres dbname=pvzhm host=localhost port=5435 sslmode=disable"
	OutboxPath   = "./db/outboxDB.json"
)
