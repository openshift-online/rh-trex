package db

import (
	"context"
	"database/sql"

	"github.com/openshift-online/rh-trex-core/db"
	"github.com/openshift-online/rh-trex-core/db/transaction"
)

// By default do no roll back transaction.
// only perform rollback if explicitly set by g2.g2.MarkForRollback(ctx, err)
const defaultRollbackPolicy = false

// newTransaction constructs a new Transaction object.
// Deprecated: Use github.com/openshift-online/rh-trex-core/db.NewTransaction instead
func newTransaction(ctx context.Context, connection SessionFactory) (*transaction.Transaction, error) {
	// Convert SessionFactory to core library interface
	coreConnection := &sessionFactoryAdapter{connection}
	return db.NewTransaction(ctx, coreConnection)
}

// sessionFactoryAdapter adapts TRex SessionFactory to core library interface
type sessionFactoryAdapter struct {
	sf SessionFactory
}

func (s *sessionFactoryAdapter) DirectDB() db.CoreDirectConnection {
	return &directConnectionAdapter{s.sf.DirectDB()}
}

// directConnectionAdapter adapts *sql.DB to CoreDirectConnection
type directConnectionAdapter struct {
	sqlDB *sql.DB
}

func (d *directConnectionAdapter) Begin() (*sql.Tx, error) {
	return d.sqlDB.Begin()
}

func (d *directConnectionAdapter) QueryRow(query string, args ...interface{}) db.CoreRow {
	return &rowAdapter{d.sqlDB.QueryRow(query, args...)}
}

// rowAdapter adapts *sql.Row to CoreRow  
type rowAdapter struct {
	sqlRow *sql.Row
}

func (r *rowAdapter) Scan(dest ...interface{}) error {
	return r.sqlRow.Scan(dest...)
}
