package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient is the global variable to hold our Redis client instance
var RedisClient *redis.Client

// InitRedisClient initializes the Redis client connection
func InitRedisClient() {
	// redisHost := os.Getenv("REDIS_HOST")
	// redisPort := os.Getenv("REDIS_PORT")
	// redisPassword := os.Getenv("REDIS_PASSWORD")
	redisUser := os.Getenv("REDIS_USER") // Not used in go-redis, but available if needed
	if redisUser == "" {
		redisUser = "default"
	}
	//var redisDB int = 0
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatalf("REDIS_URL environment variable is not set. Please provide a Redis connection URL.")
	}
	//redisTlsEnabled := os.Getenv("REDIS_TLS") == "true"

	// redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	// // Log connection attempt (without password for security)
	// log.Printf("Attempting to connect to Redis at %s (TLS: %v, DB: %d)", redisAddr, redisTlsEnabled, redisDB)

	// options := &redis.Options{
	// 	Addr:        redisAddr,
	// 	Username:    redisUser,
	// 	Password:    redisPassword,
	// 	DialTimeout: 10 * time.Second, // Corresponds to socket connectTimeout: 10000ms
	// }

	// if redisTlsEnabled {
	// 	options.TLSConfig = &tls.Config{
	// 		MinVersion: tls.VersionTLS12,
	// 	}
	// 	// For 'rediss://' scheme, go-redis will internally handle some TLS configuration.
	// 	// If using a URL, you could also parse it directly:
	// 	// parsedURL, _ := url.Parse(fmt.Sprintf("rediss://%s:%s@%s", redisUser, redisPassword, redisAddr))
	// 	// options, _ = redis.ParseURL(parsedURL.String())
	// 	encodedPassword := url.QueryEscape(redisPassword)
	// 	parsedURL, _ := url.Parse("rediss://default:" + encodedPassword + "@" + redisAddr)
	// 	options, _ = redis.ParseURL(parsedURL.String())
	// 	log.Printf("Attempting to connect to Redis with TLS via URL: %s", parsedURL)
	// }
	log.Printf("Attempting to connect to Redis using URL: %s", redisURL)

	// Parse the Redis URL into redis.Options
	// redis.ParseURL handles all parts of the URL including host, port, password, user, DB, and TLS.
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse REDIS_URL '%s': %v", redisURL, err)
	}
	// Apply a reasonable dial timeout. ParseURL might not set this.
	if options.DialTimeout == 0 {
		options.DialTimeout = 10 * time.Second
	}

	RedisClient = redis.NewClient(options)
	// Ping the Redis server to verify the connection
	// go-redis client connects lazily, so Ping() forces a connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Add a timeout for the ping
	defer cancel()
	statusCmd := RedisClient.Ping(ctx)
	if statusCmd.Err() != nil {
		log.Fatalf("Failed to connect to Redis: %v", statusCmd.Err())
	}

	log.Println("Successfully connected to Redis!")
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
