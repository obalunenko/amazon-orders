package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/obalunenko/amazon-orders/moneyutils"
)

// csv header:
// Website,"Order ID","Order Date","Purchase Order Number","Currency","Unit Price","Unit Price Tax","Shipping Charge",
// "Total Discounts","Total Owed","Shipment Item Subtotal","Shipment Item Subtotal Tax","ASIN","Product Condition",
// "Quantity","Payment Instrument Type","Order Status","Shipment Status","Ship Date","Shipping Option",
// "Shipping Address","Billing Address","Carrier Name & Tracking Number","Product Name","Gift Message",
// "Gift Sender Name","Gift Recipient Contact Details","Item Serial Number"

type order struct {
	Website                    string
	OrderID                    string
	OrderDate                  string // 2024-11-03T12:32:17Z layout
	PurchaseOrderNumber        string
	Currency                   string
	UnitPrice                  float64
	UnitPriceTax               float64
	ShippingCharge             float64
	TotalDiscounts             float64
	TotalOwed                  float64
	ShipmentItemSubtotal       string
	ShipmentItemSubtotalTax    string
	ASIN                       string
	ProductCondition           string
	Quantity                   int
	PaymentInstrumentType      string
	OrderStatus                orderStatus
	ShipmentStatus             shipmentStatus
	ShipDate                   string
	ShippingOption             string
	ShippingAddress            string
	BillingAddress             string
	CarrierNameWTrackingNumber string
	ProductName                string
	GiftMessage                string
	GiftSenderName             string
	GiftRecipientContactDetail string
	ItemSerialNumber           string
}

type col uint // column number

const (
	ColWebsite col = iota
	ColOrderID
	ColOrderDate
	ColPurchaseOrderNumber
	ColCurrency
	ColUnitPrice
	ColUnitPriceTax
	ColShippingCharge
	ColTotalDiscounts
	ColTotalOwed
	ColShipmentItemSubtotal
	ColShipmentItemSubtotalTax
	ColASIN
	ColProductCondition
	ColQuantity
	ColPaymentInstrumentType
	ColOrderStatus
	ColShipmentStatus
	ColShipDate
	ColShippingOption
	ColShippingAddress
	ColBillingAddress
	ColCarrierNameWTrackingNumber
	ColProductName
	ColGiftMessage
	ColGiftSenderName
	ColGiftRecipientContactDetail
	ColItemSerialNumber
)

// Website,"Order ID","Order Date","Purchase Order Number","Currency","Unit Price","Unit Price Tax","Shipping Charge",
// "Total Discounts","Total Owed","Shipment Item Subtotal","Shipment Item Subtotal Tax","ASIN","Product Condition",
// "Quantity","Payment Instrument Type","Order Status","Shipment Status","Ship Date","Shipping Option",
// "Shipping Address","Billing Address","Carrier Name & Tracking Number","Product Name","Gift Message",
// "Gift Sender Name","Gift Recipient Contact Details","Item Serial Number"
var colNames = map[col]string{
	ColWebsite:                    "Website",
	ColOrderID:                    "Order ID",
	ColOrderDate:                  "Order Date",
	ColPurchaseOrderNumber:        "Purchase Order Number",
	ColCurrency:                   "Currency",
	ColUnitPrice:                  "Unit Price",
	ColUnitPriceTax:               "Unit Price Tax",
	ColShippingCharge:             "Shipping Charge",
	ColTotalDiscounts:             "Total Discounts",
	ColTotalOwed:                  "Total Owed",
	ColShipmentItemSubtotal:       "Shipment Item Subtotal",
	ColShipmentItemSubtotalTax:    "Shipment Item Subtotal Tax",
	ColASIN:                       "ASIN",
	ColProductCondition:           "Product Condition",
	ColQuantity:                   "Quantity",
	ColPaymentInstrumentType:      "Payment Instrument Type",
	ColOrderStatus:                "Order Status",
	ColShipmentStatus:             "Shipment Status",
	ColShipDate:                   "Ship Date",
	ColShippingOption:             "Shipping Option",
	ColShippingAddress:            "Shipping Address",
	ColBillingAddress:             "Billing Address",
	ColCarrierNameWTrackingNumber: "Carrier Name & Tracking Number",
	ColProductName:                "Product Name",
	ColGiftMessage:                "Gift Message",
	ColGiftSenderName:             "Gift Sender Name",
	ColGiftRecipientContactDetail: "Gift Recipient Contact Details",
	ColItemSerialNumber:           "Item Serial Number",
}

func (c col) String() string {
	if name, ok := colNames[c]; ok {
		return name
	}

	return "Unknown"
}

const notAvailable = "Not Available"

type orderStatus string

func (s orderStatus) Equals(o orderStatus) bool {
	return strings.EqualFold(s.String(), o.String())
}

func (s orderStatus) String() string {
	return string(s)
}

const (
	orderStatusUnknown  orderStatus = "Unknown" // default value
	orderStatusClosed   orderStatus = "Closed"
	orderStatusCanceled orderStatus = "Cancelled"
)

func parseOrderStatus(s string) orderStatus {
	switch strings.ToLower(s) {
	case strings.ToLower(orderStatusClosed.String()):
		return orderStatusClosed
	case strings.ToLower(orderStatusCanceled.String()):
		return orderStatusCanceled
	default:
		return orderStatusUnknown
	}
}

type shipmentStatus string

func (s shipmentStatus) String() string {
	return string(s)
}

const (
	shipmentStatusUnknown      shipmentStatus = "Unknown" // default value
	shipmentStatusShipped      shipmentStatus = "Shipped"
	shipmentStatusNotAvailable shipmentStatus = notAvailable
)

func parseShipmentStatus(s string) shipmentStatus {
	switch strings.ToLower(s) {
	case strings.ToLower(shipmentStatusShipped.String()):
		return shipmentStatusShipped
	case strings.ToLower(shipmentStatusNotAvailable.String()):
		return shipmentStatusNotAvailable
	default:
		return shipmentStatusUnknown
	}
}

func calculateSpends(orders []order) float64 {
	// Calculate the total spend sum of all orders
	// Sum of all orders = sum(unit_price * quantity)
	var sum float64

	for i, o := range orders {
		if o.OrderStatus.Equals(orderStatusCanceled) {
			log.Printf("[WARN]: [order:%d] order is canceled, skipping: %s\n", i+1, o.OrderID)

			continue
		}

		if o.Quantity == 0 {
			log.Printf("[WARN]: [order:%d] quantity is 0, assuming 1: %s\n", i+1, o.OrderID)
			// Make assumption that quantity is 1 if it is not provided.
			o.Quantity = 1
		}

		sum += moneyutils.Add(moneyutils.Multiply(o.UnitPrice, float64(o.Quantity)), o.ShippingCharge)
	}

	return sum
}

func parseCSV(f *os.File) ([]order, error) {
	csvReader := csv.NewReader(f)
	csvReader.Comma = ','
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = 28

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv: %v", err)
	}

	orders := make([]order, 0, len(records)-1) // skip header.

	for i, record := range records {
		if i == 0 {
			continue
		}

		unitPrice, err := parsePrice(record[ColUnitPrice])
		if err != nil {
			return nil, fmt.Errorf("[line: %d]: failed to parse unit price: %v", i+1, err)
		}

		unitPriceTax, err := parsePrice(record[ColUnitPriceTax])
		if err != nil {
			return nil, fmt.Errorf("[line: %d]: failed to parse unit price tax: %v", i+1, err)
		}

		shippingCharge, err := parsePrice(record[ColShippingCharge])
		if err != nil {
			return nil, fmt.Errorf("[line: %d]: failed to parse shipping charge: %v", i+1, err)
		}

		totalDiscounts, err := parsePrice(record[ColTotalDiscounts])
		if err != nil {
			return nil, fmt.Errorf("[line: %d]: failed to parse total discounts: %v", i+1, err)
		}

		totalOwed, err := parsePrice(record[ColTotalOwed])
		if err != nil {
			return nil, fmt.Errorf("[line: %d]: failed to parse total owed: %v", i+1, err)
		}

		quantity, err := strconv.Atoi(record[ColQuantity])
		if err != nil {
			return nil, fmt.Errorf("[line: %d]: failed to parse quantity: %v", i+1, err)
		}

		o := order{
			Website:                    record[ColWebsite],
			OrderID:                    record[ColOrderID],
			OrderDate:                  record[ColOrderDate],
			PurchaseOrderNumber:        record[ColPurchaseOrderNumber],
			Currency:                   record[ColCurrency],
			UnitPrice:                  unitPrice,
			UnitPriceTax:               unitPriceTax,
			ShippingCharge:             shippingCharge,
			TotalDiscounts:             totalDiscounts,
			TotalOwed:                  totalOwed,
			ShipmentItemSubtotal:       record[ColShipmentItemSubtotal],
			ShipmentItemSubtotalTax:    record[ColShipmentItemSubtotalTax],
			ASIN:                       record[ColASIN],
			ProductCondition:           record[ColProductCondition],
			Quantity:                   quantity,
			PaymentInstrumentType:      record[ColPaymentInstrumentType],
			OrderStatus:                parseOrderStatus(record[ColOrderStatus]),
			ShipmentStatus:             parseShipmentStatus(record[ColShipmentStatus]),
			ShipDate:                   record[ColShipDate],
			ShippingOption:             record[ColShippingOption],
			ShippingAddress:            record[ColShippingAddress],
			BillingAddress:             record[ColBillingAddress],
			CarrierNameWTrackingNumber: record[ColCarrierNameWTrackingNumber],
			ProductName:                record[ColProductName],
			GiftMessage:                record[ColGiftMessage],
			GiftSenderName:             record[ColGiftSenderName],
			GiftRecipientContactDetail: record[ColGiftRecipientContactDetail],
			ItemSerialNumber:           record[ColItemSerialNumber],
		}

		orders = append(orders, o)
	}

	return orders, nil
}

var r = regexp.MustCompile("[$,_]")

func parsePrice(s string) (float64, error) {
	if s == "" || strings.EqualFold(s, notAvailable) {
		log.Println("[WARN]: price is not available")

		return 0, nil
	}

	s = strings.Trim(s, "'")

	return moneyutils.ParseFormatted(s, r)
}

func getOrdersCurrencies(orders []order) []string {
	if len(orders) == 0 {
		return nil
	}

	var (
		seen       = make(map[string]struct{})
		currencies = make([]string, 0, 1) // Make assumption that all orders have the same currency.
	)

	for _, o := range orders {
		if _, ok := seen[o.Currency]; !ok {
			seen[o.Currency] = struct{}{}
			currencies = append(currencies, o.Currency)
		}
	}

	return currencies
}

func getOrdersByCurrency(orders []order, cur string) []order {
	result := make([]order, 0, len(orders))

	for _, o := range orders {
		if strings.EqualFold(o.Currency, cur) {
			result = append(result, o)
		}
	}

	return result
}
