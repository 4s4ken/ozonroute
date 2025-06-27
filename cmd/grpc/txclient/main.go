package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/config"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/storage"
)

func main() {
	ctx := context.Background()

	writePool, err := pgxpool.New(ctx, config.DBConnString)
	if err != nil {
		log.Fatal("failed to connect write db:", err)
	}
	defer writePool.Close()

	readPool, err := pgxpool.New(ctx, config.DBConnString)
	if err != nil {
		log.Fatal("failed to connect read db:", err)
	}
	defer readPool.Close()

	client := storage.New(writePool, readPool)

	exec := client.Get(storage.ModeRead)
	row := exec.QueryRow(ctx, "SELECT count(*) FROM orders")
	var before int
	if err := row.Scan(&before); err != nil {
		log.Fatal("query before failed:", err)
	}
	log.Printf("Orders before: %d", before)

	err = client.Tx(ctx, storage.ModeWrite, func(tx storage.DbManager) error {
		query := `
			INSERT INTO orders 
			(order_id, user_id, status, expires_at, weight, price, package_type)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			
		`

		_, err := tx.Exec(ctx, query,
			11,
			11,
			1,
			time.Now().Add(48*time.Hour),
			123,
			123,
			2,
		)
		return err

	})
	if err != nil {
		log.Fatal("insert error:", err)
	}

	row = exec.QueryRow(ctx, "SELECT count(*) FROM orders")
	var after int
	if err := row.Scan(&after); err != nil {
		log.Fatal("query after failed:", err)
	}
	log.Printf("Orders after: %d", after)
}
