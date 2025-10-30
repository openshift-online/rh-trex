package mocks

import (
	"context"
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/openshift-online/rh-trex/pkg/db"
)

var _ db.SessionFactory = &MockSessionFactory{}

type MockSessionFactory struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
	mock   sqlmock.Sqlmock
}

// NewMockSessionFactory creates a SessionFactory using go-sqlmock.
// This provides a mock database without requiring PostgreSQL or SQLite.
func NewMockSessionFactory() *MockSessionFactory {
	// Create mock SQL database
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic("failed to create sqlmock: " + err.Error())
	}

	// Open GORM with the mock database using postgres dialector
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to create gorm DB with sqlmock: " + err.Error())
	}

	return &MockSessionFactory{
		gormDB: gormDB,
		sqlDB:  sqlDB,
		mock:   mock,
	}
}

func (m *MockSessionFactory) Init(config *config.DatabaseConfig) {
	// Mock implementation - does nothing
}

func (m *MockSessionFactory) DirectDB() *sql.DB {
	return m.sqlDB
}

func (m *MockSessionFactory) New(ctx context.Context) *gorm.DB {
	return m.gormDB.WithContext(ctx)
}

func (m *MockSessionFactory) CheckConnection() error {
	return nil
}

func (m *MockSessionFactory) Close() error {
	if m.sqlDB != nil {
		return m.sqlDB.Close()
	}
	return nil
}

func (m *MockSessionFactory) ResetDB() {
	// Mock implementation - does nothing
}

func (m *MockSessionFactory) NewListener(ctx context.Context, channel string, callback func(id string)) {
	// Mock implementation - does nothing
}
