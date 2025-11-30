package app

import (
	"confam-api/internal/config"
	"confam-api/internal/database"
	"confam-api/internal/redis"
	"confam-api/internal/server"
	"confam-api/internal/validate"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// Application holds the dependencies for the HTTP server
type Application struct {
	Config *config.Config
	Redis  *redis.Client
	Router *gin.Engine
}

// New initializes the application and its dependencies
func New(ctx context.Context) (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := database.Connect(cfg.DB); err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}
	fmt.Println("Database connected.")

	validate.SetupValidator()

	// 4. Connect to Redis
	rdb, err := redis.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}
	fmt.Println("Redis connected.")

	router := server.NewRouter(rdb)

	// if err := seeders.Seed(
	// 	database.DB,
	// 	rdb,
	// ); err != nil {
	// 	log.Fatal(err)
	// }

	return &Application{
		Config: cfg,
		Redis:  rdb,
		Router: router,
	}, nil
}

// Run starts the HTTP server
func (a *Application) Run() error {
	log.Printf("Server starting on port :%s", a.Config.Port)
	return a.Router.Run(":" + a.Config.Port)
}

func (a *Application) Close() {
	if a.Redis != nil {
		if err := a.Redis.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		} else {
			log.Println("Redis connection closed.")
		}
	}
}
