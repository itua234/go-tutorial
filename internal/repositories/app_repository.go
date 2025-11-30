package repositories

import (
	"confam-api/internal/models"
	client "confam-api/internal/redis"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type IAppRepository interface {
	CreateApp(ctx context.Context, app *models.App) error
}

type AppRepository struct {
	db    *gorm.DB
	Redis *client.Client
}

func NewAppRepository(db *gorm.DB, rdb *client.Client) *AppRepository {
	return &AppRepository{
		db:    db,
		Redis: rdb,
	}
}

func (r *AppRepository) CreateApp(ctx context.Context, app *models.App) error {
	if err := r.db.WithContext(ctx).Create(app).Error; err != nil {
		return err
	}

	testSecret := app.TestSecretKey()
	liveSecret := app.LiveSecretKey()
	fmt.Println("Test Secret:", testSecret)
	fmt.Println("Live Secret:", liveSecret)

	// Write Test Secret to Redis
	if err := r.Redis.Set(
		ctx,
		fmt.Sprintf("secret:%s", testSecret),
		app.ID,
		0,
	).Err(); err != nil {
		return fmt.Errorf("failed to save test secret to redis: %w", err)
	}

	// Write Live Secret to Redis
	if err := r.Redis.Set(
		ctx,
		fmt.Sprintf("secret:%s", liveSecret),
		app.ID,
		0,
	).Err(); err != nil {
		return fmt.Errorf("failed to save live secret to redis: %w", err)
	}

	return nil
}
