package statistics_repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (r *Repository) DeleteAllByUserID(ctx context.Context, userID models.ID) error {
	query, args := r.newDeleteAllByUserIDQuery(userID)
	_, err := r.db.Exec(ctx, query, args)

	return err
}

func (r *Repository) newDeleteAllByUserIDQuery(userID models.ID) (query string, args pgx.NamedArgs) {
	query = `
		update statistics
		set is_deleted = true
		where user_id = @user_id
	`

	args = pgx.NamedArgs{
		"user_id": int32(userID),
	}

	return
}
