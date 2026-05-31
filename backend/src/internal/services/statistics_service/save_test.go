package statistics_service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
	"github.com/ruslanonly/blindtyping/src/internal/services/statistics_service/mocks"
)

func mustParseTestTime(timeString string) time.Time {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		panic(err)
	}

	return t
}

var testUser = models.User{
	ID:        1,
	Email:     "blindtyping@gmail.com",
	Nickname:  "ffh255",
	Provider:  "google",
	CreatedAt: mustParseTestTime("2025-07-19T19:02:29+03:00"),
}

var testSaveIn = SaveIn{
	UserID:                     testUser.ID,
	WPM:                        100,
	CPM:                        300,
	Accuracy:                   99,
	Duration:                   time.Second * 25,
	Language:                   "english",
	Mode:                       "words",
	SubMode:                    "25",
	IsPunctuation:              false,
	UncompletedTestsDurationMs: 0,
	UncompletedTestsCount:      0,
	UID:                        "2e4c0f26-41ac-4750-9af6-f495b1ce71e8",
	Sign:                       "2f50408760a98cca8a713f84bd005e06",
	CreatedAt:                  mustParseTestTime("2025-10-19T19:02:29+03:00"),
	StartedAt:                  mustParseTestTime("2025-10-19T19:00:20+03:00"),
	FinishedAt:                 mustParseTestTime("2025-10-19T19:02:28+03:00"),
}

var testSignPayload = models.SignPayload{
	UID:                      testSaveIn.UID,
	WPM:                      testSaveIn.WPM,
	CPM:                      testSaveIn.CPM,
	Accuracy:                 testSaveIn.Accuracy,
	Duration:                 testSaveIn.Duration,
	Language:                 testSaveIn.Language,
	Mode:                     testSaveIn.Mode,
	SubMode:                  testSaveIn.SubMode,
	IsPunctuation:            testSaveIn.IsPunctuation,
	UncompletedTestsDuration: testSaveIn.UncompletedTestsDurationMs,
	UncompletedTestsCount:    testSaveIn.UncompletedTestsCount,
	CreatedAt:                testSaveIn.CreatedAt,
	StartedAt:                testSaveIn.StartedAt,
	FinishedAt:               testSaveIn.FinishedAt,
}

var testLanguages = models.Languages{
	"english": {},
	"russian": {},
}

func TestService_Save(t *testing.T) {
	tests := []struct {
		name   string
		in     *SaveIn
		expect func(
			statisticsRepository *mocks.MockstatisticsRepository,
			userRepository *mocks.MockuserRepository,
			languageRepository *mocks.MocklanguageRepository,
			antifroad *mocks.Mockantifroad,
			pbService *mocks.MockpbService,
			leaderboardService *mocks.MockleaderboardService,
		)
		want    *SaveOut
		wantErr bool
	}{
		{
			name: "success",
			in:   &testSaveIn,
			expect: func(
				statisticsRepository *mocks.MockstatisticsRepository,
				userRepository *mocks.MockuserRepository,
				languageRepository *mocks.MocklanguageRepository,
				antifroad *mocks.Mockantifroad,
				pbService *mocks.MockpbService,
				leaderboardService *mocks.MockleaderboardService,
			) {
				userRepository.EXPECT().GetOne(gomock.Any(), &user_repository.GetOneIn{
					UserID: &testSaveIn.UserID,
				}).Return(&testUser, nil)

				statisticsRepository.EXPECT().Exists(gomock.Any(), testSaveIn.UID).Return(false)

				antifroad.EXPECT().IsFroad(testSignPayload, testSaveIn.Sign).Return(false)

				languageRepository.EXPECT().Get().Return(testLanguages)

				pbService.EXPECT().IsPB(gomock.Any(), gomock.Any()).Return(true, 5.0)

				statisticsRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return(&models.Statistics{ID: 1}, nil)

				leaderboardService.EXPECT().UpdateRank(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			want: &SaveOut{
				StatisticsID: 1,
				IsPB:         true,
				WPMShift:     5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			statisticsRepository := mocks.NewMockstatisticsRepository(ctrl)
			userRepository := mocks.NewMockuserRepository(ctrl)
			languageRepository := mocks.NewMocklanguageRepository(ctrl)
			antifroad := mocks.NewMockantifroad(ctrl)
			pbService := mocks.NewMockpbService(ctrl)
			leaderboardService := mocks.NewMockleaderboardService(ctrl)

			tt.expect(
				statisticsRepository,
				userRepository,
				languageRepository,
				antifroad,
				pbService,
				leaderboardService,
			)

			s := &Service{
				statisticsRepository: statisticsRepository,
				userRepository:       userRepository,
				languageRepository:   languageRepository,
				antifroad:            antifroad,
				pbService:            pbService,
				leaderboardService:   leaderboardService,
			}

			got, err := s.Save(context.Background(), tt.in)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
