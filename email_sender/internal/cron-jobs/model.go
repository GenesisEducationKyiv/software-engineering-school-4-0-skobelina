package cronjobs

type Subscriber struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Event struct {
	EventID     string    `json:"eventId"`
	EventType   string    `json:"eventType"`
	AggregateID string    `json:"aggregateId"`
	Timestamp   string    `json:"timestamp"`
	Data        EventData `json:"data"`
}

type EventData struct {
	CreatedAt    string `json:"createdAt"`
	ExchangeRate string `json:"exchangeRate"`
	Email        string `json:"email"`
}
