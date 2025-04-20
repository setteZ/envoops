package utils

import (
	"os"
)

// Gets default value passed if no value exist for given environment variable.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
