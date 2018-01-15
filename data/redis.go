package data

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/subosito/gotenv"
)

var (
	client        *redis.Client
	redisPassword string
	redisHostname string
	redisPort     string
)

func init() {
	gotenv.Load()

	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisHost = os.Getenv("REDIS_HOSTNAME")
	redisPort = os.Getenv("REDIS_PORT")
}

func redisConn() *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHostname, redisPort),
			Password: redisPassword,
			DB:       0,
		})
	}
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	return client
}
