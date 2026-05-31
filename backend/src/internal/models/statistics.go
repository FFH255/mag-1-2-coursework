package models

import (
	"fmt"
	"time"
)

type Language string

type Languages map[Language]struct{}

func newLanguage(value string, languages Languages) (Language, error) {
	if value == "" {
		return "", ValidationError{
			Field:   "language",
			Message: "language cannot be blank",
		}
	}

	language := Language(value)
	_, ok := languages[language]
	if !ok {
		return "", ValidationError{
			Field:   "language",
			Message: fmt.Sprintf("language '%s' is not supported", value),
		}
	}

	return language, nil
}

type Mode string

const (
	ModeWords Mode = "words"
	ModeTime  Mode = "time"
)

type Submode string

const (
	SubmodeWords10  Submode = "10"
	SubmodeWords25  Submode = "25"
	SubmodeWords50  Submode = "50"
	SubmodeWords75  Submode = "75"
	SubmodeWords100 Submode = "100"

	SubmodeTime15s Submode = "15s"
	SubmodeTime30s Submode = "30s"
	SubmodeTime45s Submode = "45s"
	SubmodeTime1m  Submode = "1m"
	SubmodeTime2m  Submode = "2m"
)

var modes = map[Mode]map[Submode]struct{}{
	ModeWords: {
		SubmodeWords10:  {},
		SubmodeWords25:  {},
		SubmodeWords50:  {},
		SubmodeWords75:  {},
		SubmodeWords100: {},
	},
	ModeTime: {
		SubmodeTime15s: {},
		SubmodeTime30s: {},
		SubmodeTime45s: {},
		SubmodeTime1m:  {},
		SubmodeTime2m:  {},
	},
}

func newMode(mode, submode string) (Mode, Submode, error) {
	if mode == "" {
		return "", "", ValidationError{
			Field:   "mode",
			Message: "mode cannot be blank",
		}
	}

	if submode == "" {
		return "", "", ValidationError{
			Field:   "submode",
			Message: "submode cannot be blank",
		}
	}

	submodes, ok := modes[Mode(mode)]
	if !ok {
		return "", "", ValidationError{
			Field:   "mode",
			Message: fmt.Sprintf("mode '%s' does not exist", mode),
		}
	}

	_, ok = submodes[Submode(submode)]
	if !ok {
		return "", "", ValidationError{
			Field:   "submode",
			Message: fmt.Sprintf("submode '%s' does not exist", submode),
		}
	}

	return Mode(mode), Submode(submode), nil
}

type WPM float64

const (
	maxWPM = 500
	minWPM = 0
)

func newWPM(value float64) (WPM, error) {
	if value < minWPM {
		return 0, ValidationError{
			Field:   "wpm",
			Message: fmt.Sprintf("wpm cannot be less than %d", minWPM),
		}
	}

	if value > maxWPM {
		return 0, ValidationError{
			Field:   "wpm",
			Message: fmt.Sprintf("wpm cannot be greater than %d", maxWPM),
		}
	}

	return WPM(value), nil
}

type CPM float64

const (
	maxCPM = 5_000
	minCPM = 0
)

func newCPM(value float64) (CPM, error) {
	if value < minCPM {
		return 0, ValidationError{
			Field:   "cpm",
			Message: fmt.Sprintf("cpm cannot be less than %d", minCPM),
		}
	}

	if value > maxCPM {
		return 0, ValidationError{
			Field:   "cpm",
			Message: fmt.Sprintf("cpm cannot be greater than %d", maxCPM),
		}
	}

	return CPM(value), nil
}

type Accuracy float64

func newAccuracy(value float64) (Accuracy, error) {
	if value < 0 {
		return 0, ValidationError{
			Field:   "accuracy",
			Message: fmt.Sprintf("accuracy cannot be less than zero"),
		}
	}

	if value > 100 {
		return 0, ValidationError{
			Field:   "accuracy",
			Message: fmt.Sprintf("accuracy cannot be greater than 100"),
		}
	}

	return Accuracy(value), nil
}

func newDuration(value time.Duration) (time.Duration, error) {
	if value < 0 {
		return 0, ValidationError{
			Field:   "duration",
			Message: fmt.Sprintf("duration cannot be less than zero"),
		}
	}

	if value > time.Hour {
		return 0, ValidationError{
			Field:   "duration",
			Message: fmt.Sprintf("duration cannot be greater than %s", time.Hour),
		}
	}

	return value, nil
}

type Statistics struct {
	ID                            ID
	UserID                        ID
	WPM                           WPM
	CPM                           CPM
	Accuracy                      Accuracy
	Duration                      time.Duration
	PlayedAt                      time.Time
	Language                      Language
	Mode                          Mode
	SubMode                       Submode
	IsPunctuation                 bool
	UncompletedTestsCount         uint64
	UncompletedTestsTotalDuration time.Duration
	IdempotencyKey                string
}

type StatisticsOptions struct {
	UserID                        ID
	WPM                           float64
	CPM                           float64
	Accuracy                      float64
	Duration                      time.Duration
	Language                      string
	Mode                          string
	Submode                       string
	IsPunctuation                 bool
	UncompletedTestsCount         uint64
	UncompletedTestsTotalDuration time.Duration
	IdempotencyKey                string
}

func NewStatistics(opts StatisticsOptions, languages Languages) (*Statistics, error) {
	language, err := newLanguage(opts.Language, languages)
	if err != nil {
		return nil, err
	}

	mode, submode, err := newMode(opts.Mode, opts.Submode)
	if err != nil {
		return nil, err
	}

	wpm, err := newWPM(opts.WPM)
	if err != nil {
		return nil, err
	}

	cpm, err := newCPM(opts.CPM)
	if err != nil {
		return nil, err
	}

	accuracy, err := newAccuracy(opts.Accuracy)
	if err != nil {
		return nil, err
	}

	duration, err := newDuration(opts.Duration)
	if err != nil {
		return nil, err
	}

	if opts.UncompletedTestsCount < 0 || opts.UncompletedTestsCount > 1000 {
		return nil, ValidationError{"uncompletedTestsCount", "should be between 0 and 1000"}
	}

	if opts.UncompletedTestsTotalDuration < 0 || opts.UncompletedTestsTotalDuration > time.Hour {
		return nil, ValidationError{"uncompletedTestsTotalDuration", "should be between 0 and 24h"}
	}

	if opts.IdempotencyKey == "" {
		return nil, ValidationError{"idempotencyKey", "idempotencyKey must be set"}
	}

	return &Statistics{
		UserID:                        opts.UserID,
		WPM:                           wpm,
		CPM:                           cpm,
		Accuracy:                      accuracy,
		Duration:                      duration,
		Language:                      language,
		Mode:                          mode,
		SubMode:                       submode,
		IsPunctuation:                 opts.IsPunctuation,
		PlayedAt:                      time.Now(),
		UncompletedTestsCount:         opts.UncompletedTestsCount,
		UncompletedTestsTotalDuration: opts.UncompletedTestsTotalDuration,
		IdempotencyKey:                opts.IdempotencyKey,
	}, nil
}
