package main

import (
	"fmt"

	"github.com/shopspring/decimal"

	"src/src/domains"
)

func main() {
	money := domains.NewMoneyObject(decimal.NewFromFloat(100.0), "USD")

	a := decimal.NewFromFloat(0.1)
	b := decimal.NewFromFloat(0.2)

	a.Sub(b)

	fmt.Println(a, money)
}
