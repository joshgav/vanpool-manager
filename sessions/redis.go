package data

import (
	"fmt"
	"log"
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
	redisHostname = os.Getenv("REDIS_HOSTNAME")
	redisPort = os.Getenv("REDIS_PORT")
}

func redisConn() (*redis.Client, error) {
	if client == nil {
		log.Printf("creating redis client for %s:%s\n", redisHostname, redisPort)
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHostname, redisPort),
			Password: redisPassword,
			DB:       0,
		})
	}
	r, err := client.Ping().Result()
	log.Printf("redis ping result: %v\n", r)
	if err != nil {
		log.Printf("failed to connect to redis: %v\n", err)
	}
	return client, err
}
