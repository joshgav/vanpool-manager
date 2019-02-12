package util

import (
	"fmt"
	"os"
	"time"
)

func GetenvOrDefault(key, _default string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = _default
	}
	return value
}

func ParseDate(date time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", date.Year(), date.Month(), date.Day())
}
