package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client = goredis.Client

func NewClient(ctx context.Context) (*Client, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUser := os.Getenv("REDIS_USER")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisTlsEnabled := os.Getenv("REDIS_TLS") == "true"
	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	options := &goredis.Options{
		Addr:        addr,
		Username:    redisUser,
		Password:    redisPassword,
		DialTimeout: 10 * time.Second,
	}
	if redisTlsEnabled {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	rdb := goredis.NewClient(options)
	// Force the connection with a Ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", options.Addr, err)
	}
	return rdb, nil
}
