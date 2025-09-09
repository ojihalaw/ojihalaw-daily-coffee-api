package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func MustParseUUID(id string) uuid.UUID {
	parsed, err := uuid.Parse(id)
	if err != nil {
		panic(fmt.Sprintf("invalid UUID: %s", id))
	}
	return parsed
}
