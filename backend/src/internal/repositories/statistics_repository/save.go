package statistics_repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (r *Repository) Save(ctx context.Context, stats *models.Statistics) (*models.Statistics, error) {
	const query = `
		INSERT INTO statistics (
			user_id, wpm, cpm, accuracy, duration, 
		    played_at, language, mode, sub_mode, is_punctuation,
			uncompleted_tests_count, uncompleted_tests_total_duration,
		    idempotency_key
		) VALUES (
			@user_id, @wpm, @cpm, @accuracy, @duration, @played_at, @language, @mode, @sub_mode, @is_punctuation,
		    @uncompleted_tests_count, @uncompleted_tests_total_duration, @idempotency_key
		)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"user_id":                          stats.UserID,
		"wpm":                              stats.WPM,
		"cpm":                              stats.CPM,
		"accuracy":                         stats.Accuracy,
		"duration":                         stats.Duration.Milliseconds(),
		"played_at":                        stats.PlayedAt.Truncate(time.Second),
		"language":                         stats.Language,
		"mode":                             stats.Mode,
		"sub_mode":                         stats.SubMode,
		"is_punctuation":                   stats.IsPunctuation,
		"uncompleted_tests_count":          stats.UncompletedTestsCount,
		"uncompleted_tests_total_duration": stats.UncompletedTestsTotalDuration.Milliseconds(),
		"idempotency_key":                  stats.IdempotencyKey,
	}

	if err := r.db.QueryRow(ctx, query, args).Scan(&stats.ID); err != nil {
		return nil, err
	}

	return stats, nil
}
