package utils

import (
	"edtech/pkg/uuid"
)

// GenerateUUID delegates to pkg/uuid.GenerateUUID for backward compatibility.
func GenerateUUID() (string, error) {
	return uuid.GenerateUUID()
}
