package readers

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type (
	Hook interface {
		Run(token string)
	}
	PostgreSQLHook struct {
		dbpool *pgxpool.Pool
	}
)

func NewPostgreSQLHook(connection string) (Hook, error) {
	dbpool, err := pgxpool.Connect(context.Background(), connection)
	if err != nil {
		return nil, err
	}
	return &PostgreSQLHook{
		dbpool: dbpool,
	}, nil
}

func (h *PostgreSQLHook) Run(token string) {
	_, err := h.dbpool.Exec(context.Background(), "INSERT INTO tokens (token) VALUES ($1)", token)
	if err != nil {
		log.Error().Err(err).Str("token", token).Msg("error adding token to psql db")
	}
}
