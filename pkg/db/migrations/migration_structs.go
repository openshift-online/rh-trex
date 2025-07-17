package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	coredb "github.com/openshift-online/rh-trex-core/db"
)

// gormigrate is a wrapper for gorm's migration functions that adds schema versioning and rollback capabilities.
// For help writing migration steps, see the gorm documentation on migrations: http://doc.gorm.io/database.html#migration

// MigrationList rules:
//
//  1. IDs are numerical timestamps that must sort ascending.
//     Use YYYYMMDDHHMM w/ 24 hour time for format
//     Example: August 21 2018 at 2:54pm would be 201808211454.
//
//  2. Include models inline with migrations to see the evolution of the object over time.
//     Using our internal type models directly in the first migration would fail in future clean installs.
//
//  3. Migrations must be backwards compatible. There are no new required fields allowed.
//     See $project_home/g2/README.md
//
// 4. Create one function in a separate file that returns your Migration. Add that single function call to this list.
var MigrationList = []*gormigrate.Migration{
	addDinosaurs(),
	addEvents(),

	// ADD MIGRATIONS HERE
}

// Model represents the base model struct. All entities will have this struct embedded.
// This is now defined in the core framework but aliased here for backwards compatibility
type Model = coredb.Model

// fkMigration represents a foreign key relationship for database migrations
// This is now defined in the core framework but aliased here for backwards compatibility
type fkMigration = coredb.FKMigration

// CreateFK creates foreign key constraints for database migrations
// This is now implemented in the core framework but aliased here for backwards compatibility
var CreateFK = coredb.CreateFK
