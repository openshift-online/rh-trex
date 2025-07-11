package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Model represents the base model struct. All entities will have this struct embedded.
// This provides common fields for all database entities in the core framework.
type Model struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// FKMigration represents a foreign key relationship for database migrations
type FKMigration struct {
	Model     string
	Dest      string
	Field     string
	Reference string
}

// CreateFK creates foreign key constraints for database migrations
// This utility function can be used across projects to establish referential integrity
func CreateFK(g2 *gorm.DB, fks ...FKMigration) error {
	var query = `ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s ON DELETE RESTRICT ON UPDATE RESTRICT;`
	var drop = `ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s;`

	for _, fk := range fks {
		name := fmt.Sprintf("fk_%s_%s", fk.Model, fk.Dest)

		g2.Exec(fmt.Sprintf(drop, fk.Model, name))
		if err := g2.Exec(fmt.Sprintf(query, fk.Model, name, fk.Field, fk.Reference)).Error; err != nil {
			return err
		}
	}
	return nil
}