package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient is the global variable to hold our Redis client instance
var RedisClient *redis.Client

func InitRedisClient() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUser := os.Getenv("REDIS_USER") // for ACL; leave blank or "default"
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisTlsEnabled := os.Getenv("REDIS_TLS") == "true"

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	options := &redis.Options{
		Addr:        addr,
		Username:    redisUser,
		Password:    redisPassword,
		DialTimeout: 10 * time.Second,
	}
	if redisTlsEnabled {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	RedisClient = redis.NewClient(options)

	// Force the connection with a Ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", options.Addr, err)
	}
	log.Println("🔌 Successfully connected to Redis!")
}

// CloseRedisClient closes the Redis client connection
func CloseRedisClient() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Printf("Error closing Redis client: %v", err)
		} else {
			log.Println("Redis client closed.")
		}
	}
}
