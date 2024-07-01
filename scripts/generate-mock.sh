#!/bin/bash

go get -d github.com/golang/mock/mockgen
go install github.com/golang/mock/mockgen

~/go/bin/mockgen -package=rates -self_package=github.com/skobelina/currency_converter/mocks/rates -source=./domains/rates/handler.go -mock_names Service=MockRateService -destination=./mocks/rates/mock.go
~/go/bin/mockgen -package=subscribers -self_package=github.com/skobelina/currency_converter/mocks/subscribers -source=./domains/subscribers/handler.go -mock_names Service=MockSubscriberService -destination=./mocks/subscribers/mock.go
~/go/bin/mockgen -package=mails -self_package=github.com/skobelina/currency_converter/mocks/mails -source=./domains/mails/service.go -mock_names Service=MockSubscriberService -destination=./mocks/mails/mock.go

go mod tidy

