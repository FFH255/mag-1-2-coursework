package statistics_repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) Exists(ctx context.Context, idempotencyKey string) bool {
	query, args := r.buildExistsQuery(idempotencyKey)

	var exists bool

	if err := r.db.QueryRow(ctx, query, args).Scan(&exists); err != nil {
		// TODO: log error

		return false
	}

	return exists
}

func (r *Repository) buildExistsQuery(idempotencyKey string) (query string, args pgx.NamedArgs) {
	query = `SELECT EXISTS (
		SELECT 1 FROM statistics WHERE idempotency_key = @idempotency_key AND is_deleted = FALSE
	)`

	args = pgx.NamedArgs{
		"idempotency_key": idempotencyKey,
	}

	return
}
