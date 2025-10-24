package db

import (
	"context"

	"gorm.io/gorm"

	"github.com/openshift-online/rh-trex/pkg/db/transaction"
	"github.com/openshift-online/rh-trex/pkg/logger"
)

// By default do no roll back transaction.
// only perform rollback if explicitly set by g2.g2.MarkForRollback(ctx, err)
const defaultRollbackPolicy = false

// newTransaction constructs a new Transaction object.
func newTransaction(ctx context.Context, connection SessionFactory) (*transaction.Transaction, error) {
	if connection == nil {
		// This happens in non-integration tests
		return nil, nil
	}

	dbx := connection.DirectDB()
	tx, err := dbx.Begin()
	if err != nil {
		return nil, err
	}

	// current transaction ID set by postgres.  these are *not* distinct across time
	// and do get reset after postgres performs "vacuuming" to reclaim used IDs.
	var txid int64
	row := tx.QueryRow("select txid_current()")
	if row != nil {
		err := row.Scan(&txid)
		if err != nil {
			return nil, err
		}
	}

	return transaction.Build(tx, txid, defaultRollbackPolicy), nil
}

// TransactionFunc is a function type that executes operations within a transaction context.
// If the function returns an error, the transaction will be rolled back.
// Otherwise, the transaction will be committed.
type TransactionFunc func(ctx context.Context) error

// WithTransaction executes the given function within a database transaction.
// This should be used by service layer methods that need transactional behavior.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func WithTransaction(ctx context.Context, sessionFactory SessionFactory, fn TransactionFunc) error {
	log := logger.NewOCMLogger(ctx)

	// Get a new GORM DB instance
	db := sessionFactory.New(ctx)

	// If db is nil (e.g., in unit tests with mocks), execute without transaction
	if db == nil {
		log.V(10).Info("No database session available, executing without transaction (test mode)")
		return fn(ctx)
	}

	// Start a transaction using GORM's built-in support
	return db.Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction stored in it
		// This allows nested DAO calls to use the same transaction
		txCtx := context.WithValue(ctx, "gormTx", tx)

		// Execute the function with the transaction context
		err := fn(txCtx)
		if err != nil {
			log.Infof("Transaction rolled back due to error: %v", err)
			return err
		}

		log.V(10).Info("Transaction committed successfully")
		return nil
	})
}
