// ExchangeRates
//
//	    Schemes: http
//	    Host: localhost:8080
//	    BasePath: /api
//	    Version: 1.0.0
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
// swagger:meta
package main

import (
	"log"

	"github.com/skobelina/currency_converter/api"
)

func main() {
	service := api.New()
	if err := service.Handle(); err != nil {
		log.Fatal(err)
	}
}
