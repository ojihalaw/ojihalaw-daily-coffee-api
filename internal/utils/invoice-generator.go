package utils

import (
	"fmt"
	"time"
)

func GenerateInvoice(counter int) string {
	date := time.Now().Format("20060102") // YYYYMMDD
	return fmt.Sprintf("INV-%s-%04d", date, counter)
}
