package utils

import (
	"fmt"
	"strings"
	"time"
)

// GenerateSKU bikin SKU dengan format XXX-YYMMDD
func GenerateSKU(productName string) string {
	// Ambil 3 huruf depan dari nama produk (uppercase, tanpa spasi)
	prefix := strings.ToUpper(strings.ReplaceAll(productName, " ", ""))
	if len(prefix) > 3 {
		prefix = prefix[:3]
	}

	// Tambahin timestamp biar unik
	timestamp := time.Now().Format("060102150405") // YYMMDDHHMMSS

	return fmt.Sprintf("%s-%s", prefix, timestamp)
}
