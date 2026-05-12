package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedis() *redis.Client {
	var client *redis.Client

	for i := 0; i < 10; i++ {
		client = redis.NewClient(&redis.Options{
			Addr: "redis:6379",
		})

		_, err := client.Ping(Ctx).Result()

		if err == nil {
			log.Println("Connected to Redis")
			return client
		}

		log.Println("Retrying Redis connection...")
		time.Sleep(3 * time.Second)
	}

	log.Fatal("failed to connect redis")

	return nil
}
