package data

import (
	"encoding/json"
	"github.com/satori/go.uuid"
)

type SessionState struct {
	Authenticated bool
	Token         string
	Username      string
	Fullname      string
	VanpoolName   string
	Picture       []byte
}

func NewSessionStateKey() string {
	return uuid.Must(uuid.NewV4()).String()
}

func putState(key string, value SessionState) error {
	err := redisConn().Put(key, json.Marshal(value), 0).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to put %v (%v): %v", key, value, err))
	} else {
		return nil
	}
}

func getState(key string) (SessionState, error) {
	bytes, err := redisConn().Get(key).Bytes()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to get %v: %v", key, err))
	}

	var state SessionState
	err = json.Unmarshal(bytes, &state)
	if err != nil {
		return state, errors.New(fmt.Sprintf("failed to get %v: %v", key, err))
	}
	return state, nil
}
