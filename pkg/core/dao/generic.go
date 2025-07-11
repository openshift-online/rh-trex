package dao

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/openshift-online/rh-trex/pkg/core/api"
	"gorm.io/gorm"
)

// BaseDAO provides generic CRUD operations for any GORM model
type BaseDAO[T any] struct {
	db        *gorm.DB
	tableName string
}

// NewBaseDAO creates a new base DAO
func NewBaseDAO[T any](db *gorm.DB) *BaseDAO[T] {
	dao := &BaseDAO[T]{
		db: db,
	}
	
	// Get table name from model
	var model T
	dao.tableName = dao.getTableName(model)
	
	return dao
}

// Get retrieves a record by ID
func (d *BaseDAO[T]) Get(ctx context.Context, id string) (*T, error) {
	var obj T
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&obj).Error
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

// Create creates a new record
func (d *BaseDAO[T]) Create(ctx context.Context, obj *T) (*T, error) {
	err := d.db.WithContext(ctx).Create(obj).Error
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// Replace updates an existing record
func (d *BaseDAO[T]) Replace(ctx context.Context, obj *T) (*T, error) {
	err := d.db.WithContext(ctx).Save(obj).Error
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// Delete removes a record by ID
func (d *BaseDAO[T]) Delete(ctx context.Context, id string) error {
	var obj T
	return d.db.WithContext(ctx).Where("id = ?", id).Delete(&obj).Error
}

// List retrieves records with pagination and filtering
func (d *BaseDAO[T]) List(ctx context.Context, query api.ListQuery) ([]T, error) {
	var items []T
	
	db := d.db.WithContext(ctx)
	
	// Apply search filter if provided
	if query.Search != "" {
		db = d.applySearch(db, query.Search)
	}
	
	// Apply ordering if provided
	if query.OrderBy != "" {
		db = db.Order(query.OrderBy)
	}
	
	// Apply pagination
	if query.Size > 0 {
		offset := (query.Page - 1) * query.Size
		db = db.Offset(offset).Limit(query.Size)
	}
	
	err := db.Find(&items).Error
	return items, err
}

// Count returns the total number of records matching the query
func (d *BaseDAO[T]) Count(ctx context.Context, query api.ListQuery) (int, error) {
	var count int64
	
	db := d.db.WithContext(ctx).Model(new(T))
	
	// Apply search filter if provided
	if query.Search != "" {
		db = d.applySearch(db, query.Search)
	}
	
	err := db.Count(&count).Error
	return int(count), err
}

// FindByIDs retrieves multiple records by their IDs
func (d *BaseDAO[T]) FindByIDs(ctx context.Context, ids []string) ([]T, error) {
	var items []T
	err := d.db.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error
	return items, err
}

// All retrieves all records (use with caution)
func (d *BaseDAO[T]) All(ctx context.Context) ([]T, error) {
	var items []T
	err := d.db.WithContext(ctx).Find(&items).Error
	return items, err
}

// applySearch applies search filtering to the query
func (d *BaseDAO[T]) applySearch(db *gorm.DB, search string) *gorm.DB {
	// This is a simplified implementation
	// In practice, you'd parse the search string and apply appropriate filters
	// For now, we'll just search in common text fields
	
	var model T
	modelType := reflect.TypeOf(model)
	
	// Get searchable fields (string fields)
	var searchFields []string
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Type.Kind() == reflect.String && field.Name != "ID" {
			// Convert to snake_case for database column
			columnName := toSnakeCase(field.Name)
			searchFields = append(searchFields, columnName)
		}
	}
	
	if len(searchFields) == 0 {
		return db
	}
	
	// Build OR condition for text search
	var conditions []string
	var args []interface{}
	
	for _, field := range searchFields {
		conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
		args = append(args, "%"+search+"%")
	}
	
	whereClause := strings.Join(conditions, " OR ")
	return db.Where(whereClause, args...)
}

// getTableName extracts the table name from the model
func (d *BaseDAO[T]) getTableName(model T) string {
	// Use reflection to get the table name
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	
	// Convert struct name to snake_case
	return toSnakeCase(modelType.Name())
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// Transaction executes a function within a database transaction
func (d *BaseDAO[T]) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn)
}

// GetDB returns the underlying GORM database instance
func (d *BaseDAO[T]) GetDB() *gorm.DB {
	return d.db
}

// WithDB returns a new DAO instance with a different database connection
func (d *BaseDAO[T]) WithDB(db *gorm.DB) *BaseDAO[T] {
	return &BaseDAO[T]{
		db:        db,
		tableName: d.tableName,
	}
}