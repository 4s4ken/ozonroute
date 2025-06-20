package app

import (
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/store"
	"time"
)

type Service struct {
	repo  store.OrderRepository
	nowFn func() time.Time
}

func NewService(repo store.OrderRepository, nowFn func() time.Time) *Service {
	return &Service{
		repo:  repo,
		nowFn: nowFn,
	}
}
