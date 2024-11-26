package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/obalunenko/amazon-orders/internal"
	"github.com/obalunenko/amazon-orders/internal/moneyutils"
)

var inputFile = flag.String("input", "", "input csv file")

func main() {
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("input file is required")
	}

	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	orders, err := internal.ParseCSV(f)
	if err != nil {
		log.Fatalf("failed to parse csv: %v", err)
	}

	currencies := internal.GetOrdersCurrencies(orders)
	if len(currencies) == 0 {
		log.Fatalf("no currency found")
	}

	var ordersByCurrency = make(map[string][]internal.Order)

	for _, currency := range currencies {
		ordersByCurrency[currency] = internal.GetOrdersByCurrency(orders, currency)
	}

	for currency := range ordersByCurrency {
		sum := internal.CalculateSpends(ordersByCurrency[currency])

		fmt.Printf("Total spend in %s: %s\n", currency, moneyutils.ToString(moneyutils.Round(sum, 2)))
	}
}
