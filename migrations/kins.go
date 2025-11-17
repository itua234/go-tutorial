package migrations

import (
	"confam-api/models"

	"gorm.io/gorm"
)

func createNextOfKin() *Migration {
	return &Migration{
		ID: "202509110010",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.NextOfKin{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("next_of_kins")
		},
	}
}
