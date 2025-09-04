package utils

import (
	"regexp"
	"strings"
)

func GenerateSlug(input string) string {
	// lower case
	slug := strings.ToLower(input)

	// ganti spasi dengan -
	slug = strings.ReplaceAll(slug, " ", "-")

	// hapus karakter non-alfanumerik kecuali -
	re := regexp.MustCompile(`[^a-z0-9\-]+`)
	slug = re.ReplaceAllString(slug, "")

	// hapus tanda - berulang
	re2 := regexp.MustCompile(`-+`)
	slug = re2.ReplaceAllString(slug, "-")

	// trim - di awal/akhir
	slug = strings.Trim(slug, "-")

	return slug
}
