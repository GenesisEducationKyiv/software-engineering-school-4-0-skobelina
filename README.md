# ExchangeRates v1

The first version of the ExchangeRates service.

## How to run

1. Environment variables must be set. Create `.env` file:

```bush
$ echo       
export DATABASE_URL=<DATABASE_URL>
export APP_CURRENCY_EXCHANGE_URL=<APP_CURRENCY_EXCHANGE_URL>
export APP_CURRENCY_EXCHANGE_KEY=<APP_CURRENCY_EXCHANGE_KEY>
export MAILPASS=<MAILPASS>
```

and run 

```bush
$ source .env
```

2. Run program:

`go run main.go`

## PR Pipeline

- Run `go mod tidy` to verify and update dependencies.

- Run gosec:

```
$ go install github.com/securego/gosec/v2/cmd/gosec@latest
$ gosec ./...
```

## API doc with swagger

- To generate doc follow next steps:

`go get github.com/go-swagger/go-swagger/cmd/swagger`

`go install github.com/go-swagger/go-swagger/cmd/swagger`

`mkdir output`

`~/go/bin/swagger generate spec --scan-models -o ./swagger-doc/swagger.yaml`

- To open doc:

`~/go/bin/swagger serve ./swagger-doc/swagger.yaml`

## Mock tests

### How to generate mocks

- `./scripts/generate-mock.sh`
