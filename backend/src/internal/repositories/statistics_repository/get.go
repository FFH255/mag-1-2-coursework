package statistics_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/jackc/pgx/v5"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

type SortBy string

const (
	SortByPlayedAt SortBy = "playedAt"
	SortByWPM      SortBy = "wpm"
)

type Sorting struct {
	SortBy    SortBy
	Ascending bool
}

type GetIn struct {
	UserID               models.ID
	DateFrom             *time.Time
	DateFromNotInclusive *time.Time
	Language             *models.Language
	Mode                 *models.Mode
	Submode              *models.Submode
	Sorting              *Sorting
}

func (r *Repository) Get(ctx context.Context, in *GetIn) ([]models.Statistics, error) {
	query, args, err := r.buildGetQuery(in)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var stats []models.Statistics

	for rows.Next() {
		var (
			s                             models.Statistics
			durationMs                    int64
			uncompletedTestsTotalDuration int64
		)

		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.WPM,
			&s.CPM,
			&s.Accuracy,
			&durationMs,
			&s.PlayedAt,
			&s.Language,
			&s.Mode,
			&s.SubMode,
			&s.IsPunctuation,
			&s.UncompletedTestsCount,
			&uncompletedTestsTotalDuration,
		)
		if err != nil {
			return nil, err
		}

		s.Duration = time.Duration(durationMs) * time.Millisecond
		s.UncompletedTestsTotalDuration = time.Duration(uncompletedTestsTotalDuration) * time.Millisecond
		stats = append(stats, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *Repository) buildGetQuery(in *GetIn) (query string, args pgx.NamedArgs, err error) {
	if in == nil {
		return query, args, fmt.Errorf("null GetIn")
	}

	query = `
		SELECT 
			id, user_id, wpm, cpm, accuracy, duration, played_at, 
			language, mode, sub_mode, is_punctuation,
		    uncompleted_tests_count, uncompleted_tests_total_duration
		FROM statistics
		WHERE user_id = @user_id
			AND is_deleted = FALSE
	`

	args = pgx.NamedArgs{
		"user_id": in.UserID,
	}

	if in.DateFrom != nil && in.DateFromNotInclusive != nil {
		return query, args, fmt.Errorf("date_from and date_from_not_inclusive can not be both")
	}

	if in.DateFrom != nil {
		query += " AND played_at >= @date_from"
		args["date_from"] = in.DateFrom
	}

	if in.DateFromNotInclusive != nil {
		query += " AND played_at > @date_from"
		args["date_from"] = in.DateFromNotInclusive
	}

	if pointer.Get(in.Language) != "" {
		query += " AND language = @language"
		args["language"] = *in.Language
	}

	if pointer.Get(in.Mode) != "" {
		query += " AND mode = @mode"
		args["mode"] = *in.Mode
	}

	if pointer.Get(in.Submode) != "" {
		query += " AND sub_mode = @sub_mode"
		args["sub_mode"] = *in.Submode
	}

	if in.Sorting != nil {
		switch in.Sorting.SortBy {
		case SortByPlayedAt:
			query += " ORDER BY played_at"
		case SortByWPM:
			query += " ORDER BY wpm"
		}

		if in.Sorting.Ascending {
			query += " ASC"
		} else {
			query += " DESC"
		}
	}

	return query, args, nil
}
