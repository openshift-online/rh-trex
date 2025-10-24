package mocks

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/openshift-online/rh-trex/pkg/db"
)

type MockSessionFactory struct{}

func NewMockSessionFactory() db.SessionFactory {
	return &MockSessionFactory{}
}

func (f *MockSessionFactory) Init(*config.DatabaseConfig) {
	// No-op for mock
}

func (f *MockSessionFactory) DirectDB() *sql.DB {
	return nil
}

func (f *MockSessionFactory) New(ctx context.Context) *gorm.DB {
	// Return nil since we're using DAO mocks
	return nil
}

func (f *MockSessionFactory) CheckConnection() error {
	return nil
}

func (f *MockSessionFactory) Close() error {
	return nil
}

func (f *MockSessionFactory) ResetDB() {
	// No-op for mock
}

func (f *MockSessionFactory) NewListener(ctx context.Context, channel string, callback func(id string)) {
	// No-op for mock
}
