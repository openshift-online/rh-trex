package migrations

import (
	"gorm.io/gorm"

	"github.com/go-gormigrate/gormigrate/v2"
)

func addSubscriptions() *gormigrate.Migration {
	type Subscription struct {
		Model
	}

	return &gormigrate.Migration{
		ID: "202311221526",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Subscription{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Subscription{})
		},
	}
}
