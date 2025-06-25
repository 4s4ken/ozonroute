package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func toDBPackageType(pkg *pvz.PackageType) interface{} {
	if pkg == nil {
		return nil
	}
	return string(*pkg)

}
