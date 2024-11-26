package moneyutils

import (
	"regexp"

	"github.com/shopspring/decimal"
)

// Multiply returns result of multiplication of two float64.
func Multiply(a, b float64) float64 {
	d := multiply(decimal.NewFromFloat(a), decimal.NewFromFloat(b))

	return d.InexactFloat64()
}

// Div returns result of div of two float64.
func Div(a, b float64) float64 {
	d := div(decimal.NewFromFloat(a), decimal.NewFromFloat(b))

	return d.InexactFloat64()
}

// Add returns sum of two floats.
func Add(a, b float64) float64 {
	s := add(decimal.NewFromFloat(a), decimal.NewFromFloat(b))

	return s.InexactFloat64()
}

func add(a, b decimal.Decimal) decimal.Decimal {
	return a.Add(b)
}

func div(a, b decimal.Decimal) decimal.Decimal {
	return a.Div(b)
}

func multiply(a, b decimal.Decimal) decimal.Decimal {
	return a.Mul(b)
}

// Round rounds the decimal to places decimal places.
// If places < 0, it will round the integer part to the nearest 10^(-places).
func Round(a float64, places int32) float64 {
	rounded := round(decimal.NewFromFloat(a), places)

	return rounded.InexactFloat64()
}

// Parse float from string.
func Parse(raw string) (float64, error) {
	d, err := decimal.NewFromString(raw)
	if err != nil {
		return 0, err
	}

	return d.InexactFloat64(), nil
}

// ParseFormatted returns a float from a formatted string representation.
// The second argument - replRegexp, is a regular expression that is used to find characters that should be
// removed from given decimal string representation. All matched characters will be replaced with an empty string.
//
// Example:
//
//	r := regexp.MustCompile("[$,]")
//	d1, err := ParseFormatted("$5,125.99", r)
//
//	r2 := regexp.MustCompile("[_]")
//	d2, err := ParseFormatted("1_000_000", r2)
//
//	r3 := regexp.MustCompile("[USD\\s]")
//	d3, err := ParseFormatted("5000 USD", r3)
func ParseFormatted(raw string, replRegexp *regexp.Regexp) (float64, error) {
	d, err := decimal.NewFromFormattedString(raw, replRegexp)
	if err != nil {
		return 0, err
	}

	return d.InexactFloat64(), nil
}

// ToString converts float to string.
func ToString(v float64) string {
	d := decimal.NewFromFloat(v)

	return d.String()
}

func round(amount decimal.Decimal, places int32) decimal.Decimal {
	rounded := amount.Round(places)

	return rounded
}
