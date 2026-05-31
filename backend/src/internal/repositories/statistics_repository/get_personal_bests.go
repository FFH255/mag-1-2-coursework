package statistics_repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ruslanonly/blindtyping/src/internal/models"
)

type PersonalBestRecord struct {
	UserID   models.ID
	WPM      float64
	Accuracy float64
	PlayedAt time.Time
	Language models.Language
	Mode     models.Mode
	SubMode  models.Submode
}

func (r *Repository) GetPersonalBestsPaginated(ctx context.Context, limit, offset int) ([]PersonalBestRecord, error) {
	query := `
		SELECT DISTINCT ON (s.user_id, s.language, s.mode, s.sub_mode)
			s.user_id, s.wpm, s.accuracy, s.played_at, s.language, s.mode, s.sub_mode
		FROM statistics s
		WHERE s.is_deleted = FALSE
		ORDER BY s.user_id, s.language, s.mode, s.sub_mode, s.wpm DESC, s.accuracy DESC, s.played_at ASC
		LIMIT @limit OFFSET @offset
	`

	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []PersonalBestRecord
	for rows.Next() {
		var rec PersonalBestRecord
		err := rows.Scan(
			&rec.UserID,
			&rec.WPM,
			&rec.Accuracy,
			&rec.PlayedAt,
			&rec.Language,
			&rec.Mode,
			&rec.SubMode,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
