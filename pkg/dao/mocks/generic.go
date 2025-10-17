package mocks

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/dao"
)

var _ dao.GenericDao = &genericDaoMock{}

type genericDaoMock struct {
	preload  string
	orderBy  string
	joins    string
	group    string
	wheres   []dao.Where
	model    interface{}
}

func NewGenericDao() *genericDaoMock {
	return &genericDaoMock{
		wheres: []dao.Where{},
	}
}

func (g *genericDaoMock) Fetch(offset int, limit int, resourceList interface{}) error {
	// Mock implementation - does nothing but returns no error
	return nil
}

func (g *genericDaoMock) GetInstanceDao(ctx context.Context, model interface{}) dao.GenericDao {
	return &genericDaoMock{
		model:  model,
		wheres: []dao.Where{},
	}
}

func (g *genericDaoMock) Preload(preload string) {
	g.preload = preload
}

func (g *genericDaoMock) OrderBy(orderBy string) {
	g.orderBy = orderBy
}

func (g *genericDaoMock) Joins(sql string) {
	g.joins = sql
}

func (g *genericDaoMock) Group(sql string) {
	g.group = sql
}

func (g *genericDaoMock) Where(where dao.Where) {
	g.wheres = append(g.wheres, where)
}

func (g *genericDaoMock) Count(model interface{}, total *int64) {
	// Mock implementation - sets count to 0
	*total = 0
}

func (g *genericDaoMock) Validate(resourceList interface{}) error {
	// Mock implementation - returns no error
	return nil
}

func (g *genericDaoMock) GetTableName() string {
	// Mock implementation - returns empty string
	return ""
}

func (g *genericDaoMock) GetTableRelation(fieldName string) (dao.TableRelation, bool) {
	// Mock implementation - returns empty relation and false
	return dao.TableRelation{}, false
}
