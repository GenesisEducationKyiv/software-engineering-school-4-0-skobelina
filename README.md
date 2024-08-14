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

## Logging

### The application includes logging at different levels: info, warn, error, and debug.

- Examples:
```info```: General information about application flow, such as successful initialization.
```warn```: Potentially harmful situations that do not stop the application, such as failed retries.
```error```: Error events that might still allow the application to continue running.
```debug```: Detailed information on the flow through the system, helpful for debugging.

## Metrics
### Prometheus metrics have been added to monitor the application.

- Metrics available:
```http_requests_total```: Total number of HTTP requests received, labeled by method and endpoint.
```http_request_duration_seconds```: Histogram of response latency for HTTP requests.
```email_requests_total```: Total number of email sender, labeled by status.
```email_request_duration_seconds```: Histogram of response latency for email sender.

- Metrics can be accessed at:
Main service: `http://localhost:8080/metrics`
Email sender service: `http://localhost:8081/metrics`

## Alerts Configuration

### Logging Alerts

1. **Error Alerts**:
   - **Trigger**: When any error level log is generated.
   - **Reason**: To immediately notify the team about critical issues that might affect the functionality of the application.

2. **Warning Alerts**:
   - **Trigger**: When any warning level log is generated.
   - **Reason**: To notify the team about potential issues that might need attention before they escalate.

3. **Info Alerts**:
   - **Trigger**: When an important informational log is generated, such as successful startup or shutdown.
   - **Reason**: To keep track of important events in the application's lifecycle. ("Starting API service", "Starting email sender service")

### Metrics Alerts

1. **High Request Latency**:
   - **Trigger**: When the request latency exceeds 0.5 seconds.
   - **Reason**: High latency can indicate performance issues that need to be addressed to ensure a smooth user experience.

2. **High Error Rate**:
   - **Trigger**: When the error rate exceeds 5% of total requests.
   - **Reason**: A high error rate can indicate underlying issues that need immediate attention.

3. **Database Query Time**:
   - **Trigger**: When database query times exceed 1 second.
   - **Reason**: Slow queries can affect the overall performance of the application.

4. **RabbitMQ Message Publish Failures**:
   - **Trigger**: When the number of failed message publications exceeds 1% of total messages.
   - **Reason**: Ensuring reliable message delivery is crucial for the application's communication.

### Prometheus Alerts

1. **Service Availability**:
   - **Trigger**: When the service is down or not responding for 5 minutes.
   - **Reason**: To ensure the service is always available and to quickly address downtime.

2. **CPU/Memory Usage**:
   - **Trigger**: When CPU usage exceeds 80% for 5 minutes.
   - **Reason**: To monitor and manage resource utilization effectively.

3. **Custom Application Metrics**:
   - **Trigger**: When the average subscriber creation time exceeds 2 seconds.
   - **Reason**: To ensure the application is performing optimally and to identify potential bottlenecks.
