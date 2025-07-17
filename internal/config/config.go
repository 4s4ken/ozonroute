package config

import (
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/cache"
	"time"
)

const (
	DBPath       = "./db/ordersDB.json"
	DBConnString = "user=postgres password=postgres dbname=pvzhm host=localhost port=5435 sslmode=disable"
	OutboxPath   = "./db/outboxDB.json"
)

var CacheCfg = cache.Config{
	Capacity:  100,
	TTL:       5 * time.Minute,
	TypeCache: cache.LRU,
}
