#!/bin/bash

go get -d github.com/golang/mock/mockgen
go install github.com/golang/mock/mockgen

~/go/bin/mockgen -package=rates -self_package=github.com/skobelina/currency_converter/internal/rates -source=./internal/rates/handler.go -mock_names Service=MockRateService -destination=./internal/rates/mock.go
~/go/bin/mockgen -package=subscribers -self_package=github.com/skobelina/currency_converter/internal/subscribers -source=./internal/subscribers/handler.go -mock_names Service=MockSubscriberService -destination=./internal/subscribers/mock.go
~/go/bin/mockgen -package=mails -self_package=github.com/skobelina/currency_converter/internal/mails -source=./internal/mails/service.go -mock_names Service=MockSubscriberService -destination=./internal/mails/mock.go

go mod tidy
