package data

import "testing"

func TestConnection(t *testing.T) {
	cache, err := redisConn()
	if (cache == nil) && (err == nil) {
		t.Error("both cache and err should not be nil")
	}
}
