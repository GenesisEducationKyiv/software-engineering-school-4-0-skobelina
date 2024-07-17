package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	cronjobs "github.com/skobelina/email_sender/internal/cron-jobs"
	mocks_cronjobs "github.com/skobelina/email_sender/mocks/cronjobs"
	mocks_mails "github.com/skobelina/email_sender/mocks/mails"
	mocks_queue "github.com/skobelina/email_sender/mocks/queue"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestConsumeMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueue := mocks_queue.NewMockQueue(ctrl)
	mockRepo := mocks_cronjobs.NewMockRepository(ctrl)
	mockMailService := mocks_mails.NewMockMailServiceInterface(ctrl)

	cronJobService := cronjobs.NewCronJobService(mockMailService, mockQueue, mockRepo)

	messageChannel := make(chan amqp.Delivery, 1)
	mockQueue.EXPECT().ConsumeMessages().Return(messageChannel, nil)

	subscriberList := []cronjobs.Subscriber{
		{Email: "test1@gmail.com"},
		{Email: "test2@gmail.com"},
	}

	mockRepo.EXPECT().Search().Return(subscriberList, nil)

	emailRecipients := []string{"test1@gmail.com", "test2@gmail.com"}
	mockMailService.EXPECT().SendEmail(emailRecipients, "Exchange rates notification", gomock.Any()).Return(nil)

	go cronJobService.ConsumeMessages()

	// test Subscribe
	subscribeEvent := cronjobs.Event{
		EventType: "Subscribe",
		Data: cronjobs.EventData{
			Email: "test1@gmail.com",
		},
	}
	subscribeEventData, _ := json.Marshal(subscribeEvent)
	mockRepo.EXPECT().Create(&cronjobs.Subscriber{Email: "test1@gmail.com"}).Return(nil)
	messageChannel <- amqp.Delivery{Body: subscribeEventData}

	time.Sleep(1 * time.Second)

	// test Unsubscribe
	unsubscribeEvent := cronjobs.Event{
		EventType: "Unsubscribe",
		Data: cronjobs.EventData{
			Email: "test2@gmail.com",
		},
	}
	unsubscribeEventData, _ := json.Marshal(unsubscribeEvent)
	mockRepo.EXPECT().FindByEmail("test2@gmail.com").Return(&cronjobs.Subscriber{Email: "test2@gmail.com"}, nil)
	mockRepo.EXPECT().Delete(&cronjobs.Subscriber{Email: "test2@gmail.com"}).Return(nil)
	messageChannel <- amqp.Delivery{Body: unsubscribeEventData}

	time.Sleep(1 * time.Second)

	// test CurrencyRate
	currencyRateEvent := cronjobs.Event{
		EventType: "CurrencyRate",
		Data: cronjobs.EventData{
			CreatedAt:    time.Now().Format("2006-01-02"),
			ExchangeRate: "5.50",
		},
	}
	currencyRateEventData, _ := json.Marshal(currencyRateEvent)
	messageChannel <- amqp.Delivery{Body: currencyRateEventData}

	time.Sleep(1 * time.Second)

	close(messageChannel)
	assert.True(t, true)
}
