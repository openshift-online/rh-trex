package migrations

import (
	"gorm.io/gorm"

	"github.com/go-gormigrate/gormigrate/v2"
)

func addDinosaurOrganizationIdColumn() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202406131012",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE dinosaurs ADD COLUMN IF NOT EXISTS organization_id text NULL;").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
