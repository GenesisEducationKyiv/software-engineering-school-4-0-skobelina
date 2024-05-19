package currencies

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	errors "github.com/skobelina/currency_converter/utils/errors"
)

type CurrencyData struct {
	Base string
	Date string
	Motd struct {
		Msg     string
		Url     string
		Success bool
	}
	Rates map[string]float64
}

var (
	currencyExchangeApi = os.Getenv("APP_CURRENCY_EXCHANGE_URL")
	apiKey              = os.Getenv("APP_CURRENCY_EXCHANGE_KEY")
)

func GetCurrencyRates() (float64, error) {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", currencyExchangeApi, nil)
	if err != nil {
		return 0, errors.NewInternalServerErrorf("failed to create request: %v", err)
	}
	req.Header.Add("apikey", apiKey)
	resp, err := myClient.Do(req)
	if err != nil {
		return 0, errors.NewInternalServerErrorf("failed to perform request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.NewInternalServerErrorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.NewInternalServerErrorf("failed to read response body: %v", err)
	}

	var data CurrencyData
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, errors.NewInternalServerErrorf("failed to unmarshal response: %v", err)
	}

	usdRate, usdExists := data.Rates["USD"]
	uahRate, uahExists := data.Rates["UAH"]
	if !usdExists || !uahExists {
		return 0, errors.NewItemNotFoundError("USD or UAH rate not found")
	}

	usdToEur := 1 / usdRate
	usdToUah := usdToEur * uahRate

	return usdToUah, nil
}
