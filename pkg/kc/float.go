package kc

import (
	"fmt"
	"strconv"
)

func Decimal(v float64) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", v), 64)
	return value
}
