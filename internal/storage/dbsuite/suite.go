package dbsuite

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/storage/db"
)

type Suite struct {
	suite.Suite
	postgresContainer testcontainers.Container
	connStr           string
	pool              *pgxpool.Pool
	Repo              *db.Storage
}

func (s *Suite) SetupSuite() {
	ctx := context.Background()

	err := s.initPgContainer(ctx)
	s.Require().NoError(err)

	s.pool, err = pgxpool.New(ctx, s.connStr)
	s.Require().NoError(err)

	s.Repo = db.New(s.pool)

	err = s.runMigrations(ctx)
	s.Require().NoError(err)
}

func (s *Suite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	if s.postgresContainer != nil {
		err := s.postgresContainer.Terminate(context.Background())
		s.Require().NoError(err)
	}
}

func (s *Suite) initPgContainer(ctx context.Context) error {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	var err error
	s.postgresContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            false,
	})
	if err != nil {
		return err
	}

	host, err := s.postgresContainer.Host(ctx)
	if err != nil {
		return err
	}

	port, err := s.postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return err
	}

	s.connStr = fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port())
	return nil
}

func (s *Suite) runMigrations(ctx context.Context) error {
	const createOrdersTable = `
	CREATE TABLE IF NOT EXISTS orders (
		order_id     BIGSERIAL PRIMARY KEY NOT NULL,
		user_id      BIGINT                NOT NULL,
		status       INTEGER               NOT NULL DEFAULT 0,
		weight       NUMERIC(10, 2)        NOT NULL,
		price        NUMERIC(10, 2)        NOT NULL,
		package_type VARCHAR,
		expires_at   TIMESTAMPTZ           NOT NULL,
		created_at   TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at   TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	const createHistoryTable = `
	CREATE TABLE IF NOT EXISTS order_status_history (
		order_id   BIGINT      NOT NULL,
		status     INTEGER     NOT NULL,
		changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`

	_, err := s.pool.Exec(ctx, createOrdersTable)
	if err != nil {
		return err
	}

	_, err = s.pool.Exec(ctx, createHistoryTable)
	return err
}
