#!/bin/bash

go get -d github.com/golang/mock/mockgen
go install github.com/golang/mock/mockgen

~/go/bin/mockgen -package=queue -self_package=github.com/skobelina/email_sender/mocks/queue -source=./pkg/queue/queue.go -mock_names Queue=MockQueue -destination=./mocks/queue/mock.go
~/go/bin/mockgen -package=cronjobs -self_package=github.com/skobelina/email_sender/mocks/cronjobs -source=./internal/cron-jobs/repository.go -mock_names Repository=MockRepository -destination=./mocks/cronjobs/mock.go
~/go/bin/mockgen -package=mails -self_package=github.com/skobelina/email_sender/mocks/mails -source=./internal/mails/service.go -mock_names MailService=MockMailService -destination=./mocks/mails/mock.go

go mod tidy
