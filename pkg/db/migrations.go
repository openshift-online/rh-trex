package db

import (
	"context"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/golang/glog"
	"github.com/openshift-online/rh-trex/pkg/db/migrations"

	"gorm.io/gorm"
)

// gormigrate is a wrapper for gorm's migration functions that adds schema versioning and rollback capabilities.
// For help writing migration steps, see the gorm documentation on migrations: http://doc.gorm.io/database.html#migration

func Migrate(g2 *gorm.DB) error {
	m := newGormigrate(g2)

	if err := m.Migrate(); err != nil {
		return err
	}
	return nil
}

// MigrateTo a specific migration will not seed the database, seeds are up to date with the latest
// schema based on the most recent migration
// This should be for testing purposes mainly
func MigrateTo(sessionFactory SessionFactory, migrationID string) {
	g2 := sessionFactory.New(context.Background())
	m := newGormigrate(g2)

	if err := m.MigrateTo(migrationID); err != nil {
		glog.Fatalf("Could not migrate: %v", err)
	}
}

func newGormigrate(g2 *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(g2, gormigrate.DefaultOptions, migrations.MigrationList)
}
