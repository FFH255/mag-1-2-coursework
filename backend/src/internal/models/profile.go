package models

import (
	"fmt"
	"time"
)

type PersonalBest struct {
	Language Language
	Mode     Mode
	SubMode  Submode
	WPM      WPM
	Accuracy Accuracy
	PlayedAt time.Time
}

func (pb *PersonalBest) needsUpdate(stats *Statistics) bool {
	switch {
	case stats.WPM > pb.WPM:
		return true
	case stats.WPM == pb.WPM && stats.Accuracy > pb.Accuracy:
		return true
	}

	return false
}

func (pb *PersonalBest) update(stats *Statistics) {
	pb.WPM = stats.WPM
	pb.Accuracy = stats.Accuracy
	pb.PlayedAt = stats.PlayedAt
}

func newPersonalBest(stats *Statistics) *PersonalBest {
	return &PersonalBest{
		Language: stats.Language,
		Mode:     stats.Mode,
		SubMode:  stats.SubMode,
		WPM:      stats.WPM,
		Accuracy: stats.Accuracy,
		PlayedAt: stats.PlayedAt,
	}
}

type LanguageStats struct {
	Language       Language
	TestsStarted   uint64
	TestsCompleted uint64
	TimePlayed     time.Duration
}

func newLanguageStats(language Language) *LanguageStats {
	return &LanguageStats{
		Language: language,
	}
}

func (s *LanguageStats) add(stats *Statistics) {
	if s.Language != stats.Language {
		return
	}

	s.TestsCompleted++
	s.TestsStarted++
	s.TestsStarted += stats.UncompletedTestsCount
	s.TimePlayed += stats.Duration + stats.UncompletedTestsTotalDuration
}

type personalBests map[string]*PersonalBest

func newPersonalBests() personalBests {
	return make(map[string]*PersonalBest)
}

func (pb personalBests) getPersonalBestID(stats *Statistics) string {
	return fmt.Sprintf("%s_%s_%s", stats.Language, stats.Mode, stats.SubMode)
}

func (pb personalBests) set(stats *Statistics) {
	if stats == nil {
		return
	}

	id := pb.getPersonalBestID(stats)
	personalBest, ok := pb[id]
	if !ok {
		pb[id] = newPersonalBest(stats)
		return
	}

	if !personalBest.needsUpdate(stats) {
		return
	}

	personalBest.update(stats)
}

var allowedModes = map[Mode]map[Submode]struct{}{
	ModeWords: {
		SubmodeWords10:  {},
		SubmodeWords25:  {},
		SubmodeWords50:  {},
		SubmodeWords100: {},
	},
	ModeTime: {
		SubmodeTime15s: {},
		SubmodeTime30s: {},
		SubmodeTime1m:  {},
		SubmodeTime2m:  {},
	},
}

func withAllowedMode(stats *Statistics) bool {
	submodes, ok := allowedModes[stats.Mode]
	if !ok {
		return false
	}

	_, ok = submodes[stats.SubMode]
	return ok
}

type Profile struct {
	UserID               ID
	Username             string
	JoinedAt             time.Time
	TimePlayed           time.Duration
	CompletedTests       uint64
	StartedTests         uint64
	PersonalBests        []PersonalBest
	LanguageStats        []LanguageStats
	LeaderboardPositions []LeaderboardPosition
}

func NewProfile(user *User, stats []Statistics, leaderboardPositions []LeaderboardPosition) *Profile {
	var (
		completedTests   uint64        = 0
		startedTests     uint64        = 0
		timePlayed       time.Duration = 0
		pbs                            = newPersonalBests()
		languageStatsMap               = make(map[Language]*LanguageStats)
	)

	for _, stat := range stats {
		completedTests++
		startedTests++
		startedTests += stat.UncompletedTestsCount
		timePlayed += stat.Duration + stat.UncompletedTestsTotalDuration

		if !withAllowedMode(&stat) {
			continue
		}

		if _, ok := languageStatsMap[stat.Language]; !ok {
			languageStatsMap[stat.Language] = newLanguageStats(stat.Language)
		}

		languageStats := languageStatsMap[stat.Language]
		languageStats.add(&stat)

		pbs.set(&stat)
	}

	pbList := make([]PersonalBest, 0, len(pbs))
	for _, pb := range pbs {
		pbList = append(pbList, *pb)
	}

	languageStatsList := make([]LanguageStats, 0, len(languageStatsMap))
	for _, stat := range languageStatsMap {
		languageStatsList = append(languageStatsList, *stat)
	}

	return &Profile{
		UserID:               user.ID,
		Username:             string(user.Nickname),
		JoinedAt:             user.CreatedAt,
		TimePlayed:           timePlayed,
		CompletedTests:       completedTests,
		StartedTests:         startedTests,
		PersonalBests:        pbList,
		LanguageStats:        languageStatsList,
		LeaderboardPositions: leaderboardPositions,
	}
}
