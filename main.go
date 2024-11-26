package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/obalunenko/amazon-orders/moneyutils"
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

	orders, err := parseCSV(f)
	if err != nil {
		log.Fatalf("failed to parse csv: %v", err)
	}

	currencies := getOrdersCurrencies(orders)
	if len(currencies) == 0 {
		log.Fatalf("no currency found")
	}

	var ordersByCurrency = make(map[string][]order)

	for _, currency := range currencies {
		ordersByCurrency[currency] = getOrdersByCurrency(orders, currency)
	}

	for currency := range ordersByCurrency {
		sum := calculateSpends(ordersByCurrency[currency])

		fmt.Printf("Total spend in %s: %s\n", currency, moneyutils.ToString(moneyutils.Round(sum, 2)))
	}
}
