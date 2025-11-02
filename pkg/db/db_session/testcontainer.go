package db_session

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/openshift-online/rh-trex/pkg/db"
)

type Testcontainer struct {
	config    *config.DatabaseConfig
	container *postgres.PostgresContainer
	g2        *gorm.DB
	sqlDB     *sql.DB
}

var _ db.SessionFactory = &Testcontainer{}

// NewTestcontainerFactory creates a SessionFactory using testcontainers.
// This starts a real PostgreSQL container for integration testing.
func NewTestcontainerFactory(config *config.DatabaseConfig) *Testcontainer {
	conn := &Testcontainer{
		config: config,
	}
	conn.Init(config)
	return conn
}

func (f *Testcontainer) Init(config *config.DatabaseConfig) {
	ctx := context.Background()

	glog.Infof("Starting PostgreSQL testcontainer...")

	// Create PostgreSQL container
	container, err := postgres.Run(ctx,
		"postgres:14.2",
		postgres.WithDatabase(config.Name),
		postgres.WithUsername(config.Username),
		postgres.WithPassword(config.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		glog.Fatalf("Failed to start PostgreSQL testcontainer: %s", err)
	}

	f.container = container

	// Get connection string from container
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		glog.Fatalf("Failed to get connection string from testcontainer: %s", err)
	}

	glog.Infof("PostgreSQL testcontainer started at: %s", connStr)

	// Open SQL connection
	f.sqlDB, err = sql.Open("postgres", connStr)
	if err != nil {
		glog.Fatalf("Failed to connect to testcontainer database: %s", err)
	}

	// Configure connection pool
	f.sqlDB.SetMaxOpenConns(config.MaxOpenConnections)

	// Connect GORM to use the same connection
	conf := &gorm.Config{
		PrepareStmt:            false,
		FullSaveAssociations:   false,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	}

	if config.Debug {
		conf.Logger = logger.Default.LogMode(logger.Info)
	}

	f.g2, err = gorm.Open(gormpostgres.New(gormpostgres.Config{
		Conn:                 f.sqlDB,
		PreferSimpleProtocol: true,
	}), conf)
	if err != nil {
		glog.Fatalf("Failed to connect GORM to testcontainer database: %s", err)
	}

	// Run migrations
	glog.Infof("Running database migrations on testcontainer...")
	if err := db.Migrate(f.g2); err != nil {
		glog.Fatalf("Failed to run migrations on testcontainer: %s", err)
	}

	glog.Infof("Testcontainer database initialized successfully")
}

func (f *Testcontainer) DirectDB() *sql.DB {
	return f.sqlDB
}

func (f *Testcontainer) New(ctx context.Context) *gorm.DB {
	conn := f.g2.Session(&gorm.Session{
		Context: ctx,
		Logger:  f.g2.Logger.LogMode(logger.Silent),
	})
	if f.config.Debug {
		conn = conn.Debug()
	}
	return conn
}

func (f *Testcontainer) CheckConnection() error {
	_, err := f.sqlDB.Exec("SELECT 1")
	return err
}

func (f *Testcontainer) Close() error {
	ctx := context.Background()

	// Close SQL connection
	if f.sqlDB != nil {
		if err := f.sqlDB.Close(); err != nil {
			glog.Errorf("Error closing SQL connection: %s", err)
		}
	}

	// Terminate container
	if f.container != nil {
		glog.Infof("Stopping PostgreSQL testcontainer...")
		if err := f.container.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate testcontainer: %s", err)
		}
		glog.Infof("PostgreSQL testcontainer stopped")
	}

	return nil
}

func (f *Testcontainer) ResetDB() {
	// For testcontainers, we can just truncate all tables
	ctx := context.Background()
	g2 := f.New(ctx)

	// Dynamically retrieve all table names except for the "migrations" table and truncate them
	var tableNames []string
	err := g2.Raw(`
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
		AND tablename != 'migrations'
	`).Scan(&tableNames).Error
	if err != nil {
		glog.Errorf("Error retrieving table names for reset: %s", err)
		return
	}
	for _, table := range tableNames {
		if err := g2.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			glog.Errorf("Error truncating table %s: %s", table, err)
		}
	}
}

func (f *Testcontainer) NewListener(ctx context.Context, channel string, callback func(id string)) {
	// Get the connection string for the listener
	connStr, err := f.container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		glog.Errorf("Failed to get connection string for listener: %s", err)
		return
	}

	newListener(ctx, connStr, channel, callback)
}
