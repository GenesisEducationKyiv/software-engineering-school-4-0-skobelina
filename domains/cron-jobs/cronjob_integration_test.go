package cronjobs_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	cronjobs "github.com/skobelina/currency_converter/domains/cron-jobs"
	"github.com/skobelina/currency_converter/domains/mails"
	"github.com/skobelina/currency_converter/domains/rates"
	"github.com/skobelina/currency_converter/domains/subscribers"
	"github.com/skobelina/currency_converter/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	mockDB, db, err := repo.Mock()
	require.NoError(t, err)
	return mockDB, db
}

func TestNotificationExchangeRates_Success(t *testing.T) {
	_, db := setupMockDB(t)
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("failed to get sql.DB from gorm.DB: %v", err)
		}
		sqlDB.Close()
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMailService := mails.NewMockMailServiceInterface(ctrl)
	mockRateService := rates.NewMockRateServiceInterface(ctrl)
	mockSubscriberService := subscribers.NewMockSubscriberServiceInterface(ctrl)

	mockRate := 27.32
	mockRateService.EXPECT().Get().Return(&mockRate, nil)

	mockSubscribers := &subscribers.SearchSubscribeResponse{
		Data: []subscribers.Subscriber{
			{Email: "test1@example.com"},
			{Email: "test2@example.com"},
		},
	}
	mockSubscriberService.EXPECT().Search(gomock.Any()).Return(mockSubscribers, nil)
	mockMailService.EXPECT().SendEmail(gomock.Any(), "Exchange rates notification", gomock.Any()).Return(nil)
	cronJobService := cronjobs.NewService(db, mockMailService, mockRateService, mockSubscriberService)
	err := cronJobService.NotificationExchangeRates()
	assert.NoError(t, err)
}

func TestNotificationExchangeRates_Error(t *testing.T) {
	_, db := setupMockDB(t)
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("failed to get sql.DB from gorm.DB: %v", err)
		}
		sqlDB.Close()
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMailService := mails.NewMockMailServiceInterface(ctrl)
	mockRateService := rates.NewMockRateServiceInterface(ctrl)
	mockSubscriberService := subscribers.NewMockSubscriberServiceInterface(ctrl)
	mockRateService.EXPECT().Get().Return(nil, errors.New("internal server error"))
	mockSubscriberService.EXPECT().Search(gomock.Any()).Return(&subscribers.SearchSubscribeResponse{}, nil)
	cronJobService := cronjobs.NewService(db, mockMailService, mockRateService, mockSubscriberService)
	err := cronJobService.NotificationExchangeRates()
	assert.Error(t, err)
}
