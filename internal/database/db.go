package database

import (
	"confam-api/internal/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func Connect(c Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//log.Fatal("failed to connect to database:", err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	DB = db
	return nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Request{},
		&models.Company{},
		&models.App{},
		&models.Customer{},
		&models.Identity{},
		&models.NextOfKin{},
		//&models.KycVerification{},
		//&models.BankAccount{},
	)
}
