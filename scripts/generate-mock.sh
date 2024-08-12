#!/bin/bash

go get -d github.com/golang/mock/mockgen
go install github.com/golang/mock/mockgen

~/go/bin/mockgen -package=rates -self_package=github.com/skobelina/currency_converter/mocks/rates -source=./internal/rates/handler.go -mock_names Service=MockRateService -destination=./mocks/rates/mock.go
~/go/bin/mockgen -package=subscribers -self_package=github.com/skobelina/currency_converter/mocks/subscribers -source=./internal/subscribers/handler.go -mock_names Service=MockSubscriberService -destination=./mocks/subscribers/mock.go

go mod tidy
