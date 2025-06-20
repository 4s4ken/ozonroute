package main

import (
	"context"
	"database/sql"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/config"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/storage"
)

func main() {
	writeDB, err := sql.Open("pgx", config.DBConnString)
	if err != nil {
		log.Fatal("failed to connect write db:", err)
	}
	defer writeDB.Close()

	readDB, err := sql.Open("pgx", config.DBConnString)
	if err != nil {
		log.Fatal("failed to connect read db:", err)
	}
	defer readDB.Close()

	client := storage.New(writeDB, readDB)
	ctx := context.Background()

	exec := client.Get(storage.ModeRead)
	row := exec.QueryRowContext(ctx, "SELECT count(*) FROM orders")
	var before int
	if err := row.Scan(&before); err != nil {
		log.Fatal(err)
	}
	log.Printf("Orders before: %d", before)

	err = client.Tx(ctx, func(tx storage.DbManager) error {
		query := `
		INSERT INTO orders 
		(order_id, user_id, status, expires_at, weight, price, package_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := tx.ExecContext(ctx, query,
			123,
			123,
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

	row = exec.QueryRowContext(ctx, "SELECT count(*) FROM orders")
	var after int
	if err := row.Scan(&after); err != nil {
		log.Fatal(err)
	}
	log.Printf("Orders after: %d", after)
}
