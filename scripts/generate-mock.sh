#!/bin/bash

go get -d github.com/golang/mock/mockgen
go install github.com/golang/mock/mockgen

~/go/bin/mockgen -package=rates -self_package=github.com/skobelina/currency_converter/domains/rates -source=./domains/rates/service.go -mock_names Service=MockRateService -destination=./domains/rates/mock.go
~/go/bin/mockgen -package=subscribers -self_package=github.com/skobelina/currency_converter/domains/subscribers -source=./domains/subscribers/service.go -mock_names Service=MockSubscriberService -destination=./domains/subscribers/mock.go

go mod tidy

