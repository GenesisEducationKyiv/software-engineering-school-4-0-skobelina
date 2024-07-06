package mails

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

const (
	DefaultMailSendAddress = "marisa.skobelina@gmail.com"
	DefaultMailHost        = "smtp.gmail.com"
	DefaultMailPort        = 587

	messagesBuffer             = 10
	closeSessionTimeoutSeconds = 30
	messagesLimitPerSending    = 50
	messagesCountPerMinute     = 25
)

type MailServiceInterface interface {
	SendEmail(recipients []string, subject string, temp Template) error
	SendBatch(messages ...*Message) error
}

type MailService struct {
	username string
	host     string
	messages chan []*gomail.Message
}

var _ MailServiceInterface = (*MailService)(nil)

func NewService(username string, host string) *MailService {
	s := &MailService{
		username: username,
		host:     host,

		messages: make(chan []*gomail.Message, messagesBuffer),
	}

	go s.scheduler()

	return s
}

type Message struct {
	Recipients []string
	Subject    string
	Template   Template
}

func (s *MailService) SendEmail(recipients []string, subject string, temp Template) error {
	message := &Message{
		Recipients: recipients,
		Subject:    subject,
		Template:   temp,
	}
	return s.SendBatch(message)
}

func (s *MailService) SendBatch(messages ...*Message) error {
	preparedMessages := prepare(messages...)
	for i := 0; i < len(preparedMessages); i += messagesLimitPerSending {
		first, last := i, messagesLimitPerSending+i
		if last > len(preparedMessages) {
			last = len(preparedMessages)
		}
		s.messages <- preparedMessages[first:last]
	}
	return nil
}

func prepare(messages ...*Message) []*gomail.Message {
	var gomailMessages []*gomail.Message
	for _, message := range messages {
		m := gomail.NewMessage()
		m.SetHeader("From", DefaultMailSendAddress)
		m.SetHeader("To", message.Recipients...)
		m.SetHeader("Subject", message.Subject)
		parsedMail := Parse(message.Template)
		m.SetBody("text/html", parsedMail)
		gomailMessages = append(gomailMessages, m)
	}
	return gomailMessages
}

func (s *MailService) scheduler() {
	var sendCloser gomail.SendCloser
	var err error

	d := gomail.NewDialer(DefaultMailHost, DefaultMailPort, DefaultMailSendAddress, os.Getenv("MAILPASS"))
	open := false

	for {
		select {
		case m, ok := <-s.messages:
			if !ok {
				return
			}
			if !open {
				if sendCloser, err = d.Dial(); err != nil {
					logrus.Errorf("mails: opening connection: %v", err)
					continue
				}
				open = true
			}
			if err := gomail.Send(sendCloser, m...); err != nil {
				logrus.Errorf("mails: cannot send email: %v", err)
			}
			if len(m) >= messagesCountPerMinute {
				<-time.After(time.Minute)
			}
		// Close the connection to the SMTP server if no email was sent.
		case <-time.After(closeSessionTimeoutSeconds * time.Second):
			if open {
				if err := sendCloser.Close(); err != nil {
					logrus.Errorf("mails: closing connection: %v", err)
				}
				open = false
			}
		}
	}
}
