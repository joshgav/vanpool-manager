package main

import (
	"os"
)

func GetenvOrDefault(key, _default string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = _default
	}
	return value
}
