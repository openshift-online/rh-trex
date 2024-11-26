package services

import (
	"context"
	e "errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"github.com/Masterminds/squirrel"
	"github.com/yaacov/tree-search-language/pkg/tsl"
	"github.com/yaacov/tree-search-language/pkg/walkers/ident"
	sqlFilter "github.com/yaacov/tree-search-language/pkg/walkers/sql"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/auth"
	"github.com/openshift-online/rh-trex/pkg/client/ocm"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/pkg/errors"
	"github.com/openshift-online/rh-trex/pkg/logger"
)

type GenericService interface {
	List(ctx context.Context, args *ListArguments, resourceList interface{}) (*api.PagingMeta, *errors.ServiceError)
}

func NewGenericService(genericDao dao.GenericDao, ocmClient *ocm.Client) GenericService {
	return &sqlGenericService{genericDao: genericDao, ocmClient: ocmClient}
}

var _ GenericService = &sqlGenericService{}

type sqlGenericService struct {
	genericDao dao.GenericDao
	ocmClient  *ocm.Client
}

var (
	searchDisallowedFields = map[string][]string{}
	allFieldsAllowed       = []string{}
	// Some mappings are not required as they match AMS resource 1:1
	// Such as Organization
	modelToAmsResource = map[string]string{}

	// TODO: This should be more dynamic
	// prefarably utilizing the openapi json via reflect
	// and the column names from the model
	openapiToModelFields = map[string]dao.TableMappingRelation{
		api.DinosaurTypeName: dao.DinosaurApiToModel(),
	}
)

// wrap all needed pieces for the LIST funciton
type listContext struct {
	ctx              context.Context
	args             *ListArguments
	username         string
	pagingMeta       *api.PagingMeta
	ulog             *logger.OCMLogger
	resourceList     interface{}
	disallowedFields []string
	openapiToModel   map[string]string
	resourceType     string
	joins            map[string]dao.TableRelation
	groupBy          []string
	set              map[string]bool
}

func newListContext(
	ctx context.Context,
	args *ListArguments,
	resourceList interface{},
) (*listContext, interface{}, *errors.ServiceError) {
	username := auth.GetUsernameFromContext(ctx)
	log := logger.NewOCMLogger(ctx)
	resourceModel := reflect.TypeOf(resourceList).Elem().Elem()
	resourceTypeStr := resourceModel.Name()
	if resourceTypeStr == "" {
		return nil, nil, errors.GeneralError("Could not determine resource type")
	}
	disallowedFields := searchDisallowedFields[resourceTypeStr]
	if disallowedFields == nil {
		disallowedFields = allFieldsAllowed
	}
	openapiToModel := openapiToModelFields[resourceTypeStr]
	args.Search = strings.Trim(args.Search, " ")
	return &listContext{
		ctx:              ctx,
		args:             args,
		username:         username,
		pagingMeta:       &api.PagingMeta{Page: args.Page},
		ulog:             &log,
		resourceList:     resourceList,
		disallowedFields: disallowedFields,
		openapiToModel:   openapiToModel.Mapping,
		resourceType:     resourceTypeStr,
	}, reflect.New(resourceModel).Interface(), nil
}

func resourceIncludesOrgId(model interface{}) bool {
	resourceModel := reflect.TypeOf(model).Elem()
	_, found := resourceModel.FieldByName("OrganizationId")
	return found
}

func isAllowedToAllOrgs(allowedOrgs []string) bool {
	return len(allowedOrgs) == 1 && allowedOrgs[0] == "*"
}

func (s *sqlGenericService) populateSearchRestriction(listCtx *listContext, model any) *errors.ServiceError {
	ctx := listCtx.ctx
	resourceName := listCtx.resourceType
	if name, ok := modelToAmsResource[resourceName]; ok {
		resourceName = string(name)
	}
	if resourceIncludesOrgId(model) {
		resourceReview, err := s.ocmClient.Authorization.ResourceReview(
			ctx,
			listCtx.username,
			auth.GetAction,
			resourceName,
		)
		if err != nil {
			return errors.GeneralError(
				"Failed to verify resource review for user '%s' on resource '%s': %v",
				listCtx.username,
				listCtx.resourceType,
				err,
			)
		}

		// TODO setup a search builder
		allowedOrgs := resourceReview.OrganizationIDs()
		// If user doesn't have access to all orgs include search for allowed only
		if !isAllowedToAllOrgs(allowedOrgs) {
			if listCtx.args.Search != "" {
				listCtx.args.Search += " and "
			}
			for i := range allowedOrgs {
				allowedOrgs[i] = fmt.Sprintf("'%s'", allowedOrgs[i])
			}
			listCtx.args.Search += fmt.Sprintf("organization_id in (%s)", strings.Join(allowedOrgs, ","))
		}
	}
	return nil
}

// resourceList must be a pointer to a slice of database resource objects
func (s *sqlGenericService) List(
	ctx context.Context,
	args *ListArguments,
	resourceList interface{},
) (*api.PagingMeta, *errors.ServiceError) {
	listCtx, model, err := newListContext(ctx, args, resourceList)
	if err != nil {
		return nil, err
	}

	if err = s.populateSearchRestriction(listCtx, model); err != nil {
		return nil, err
	}

	// the ordering for the sub functions matters.
	builders := []listBuilder{
		// build SQL to load related resource. for now, it delegates to gorm.preload.
		s.buildPreload,

		// add "ORDER BY"
		s.buildOrderBy,

		// translate "search" into "WHERE"(s), and "JOIN"(s) if related resource is searched.
		s.buildSearch,

		// TODO: add any custom builder functions
	}

	d := s.genericDao.GetInstanceDao(listCtx.ctx, model)

	// run all the "builders". they cumulatively add constructs to gorm by the context.
	// it stops when a builder function raises error or signals finished.
	var finished bool
	for _, builderFn := range builders {
		if finished, err = builderFn(listCtx, &d); err != nil {
			return nil, err
		}
		if finished {
			if err = s.loadList(listCtx, &d); err != nil {
				return nil, err
			}
			break
		}
	}
	return listCtx.pagingMeta, nil
}

/*** Define all sub functions in the type of listBuilder ***/
type listBuilder func(*listContext, *dao.GenericDao) (finished bool, err *errors.ServiceError)

func (s *sqlGenericService) buildPreload(listCtx *listContext, d *dao.GenericDao) (bool, *errors.ServiceError) {
	listCtx.set = make(map[string]bool)

	for _, preload := range listCtx.args.Preloads {
		listCtx.set[preload] = true
	}
	// preload each table only once; struct{} doesn't occupy any additional space
	for _, preload := range listCtx.args.Preloads {
		(*d).Preload(preload)
	}
	return false, nil
}

func (s *sqlGenericService) buildOrderBy(listCtx *listContext, d *dao.GenericDao) (bool, *errors.ServiceError) {
	if len(listCtx.args.OrderBy) != 0 {
		orderByArgs, serviceErr := db.ArgsToOrderBy(listCtx.args.OrderBy, listCtx.disallowedFields,
			listCtx.openapiToModel, (*d).GetTableName())
		if serviceErr != nil {
			return false, serviceErr
		}
		for _, orderByArg := range orderByArgs {
			(*d).OrderBy(orderByArg)
		}
	}
	return false, nil
}

func (s *sqlGenericService) buildSearchValues(
	listCtx *listContext,
	d *dao.GenericDao,
) (string, []any, *errors.ServiceError) {
	if listCtx.args.Search == "" {
		s.addJoins(listCtx, d)
		return "", nil, nil
	}

	// create the TSL tree
	tslTree, err := tsl.ParseTSL(listCtx.args.Search)
	if err != nil {
		return "", nil, errors.BadRequest("Failed to parse search query: %s", listCtx.args.Search)
	}
	// find all related tables
	tslTree, serviceErr := s.treeWalkForRelatedTables(listCtx, tslTree, d)
	if serviceErr != nil {
		return "", nil, serviceErr
	}
	// prepend table names to prevent "ambiguous" errors
	tslTree, serviceErr = s.treeWalkForAddingTableName(listCtx, tslTree, d)
	if serviceErr != nil {
		return "", nil, serviceErr
	}
	// convert to sqlizer
	_, sqlizer, serviceErr := s.treeWalkForSqlizer(listCtx, tslTree)
	if serviceErr != nil {
		return "", nil, serviceErr
	}

	s.addJoins(listCtx, d)

	// parse the search string to SQL WHERE
	sql, values, err := sqlizer.ToSql()
	if err != nil {
		return "", nil, errors.GeneralError(err.Error())
	}
	return sql, values, nil
}

func (s *sqlGenericService) buildSearch(listCtx *listContext, d *dao.GenericDao) (bool, *errors.ServiceError) {
	sql, values, err := s.buildSearchValues(listCtx, d)
	if err != nil {
		return false, err
	}
	(*d).Where(dao.NewWhere(sql, values))
	return true, nil
}

// JOIN the tables that appear in the search string
func (s *sqlGenericService) addJoins(listCtx *listContext, d *dao.GenericDao) {
	for _, r := range listCtx.joins {
		if _, ok := listCtx.set[r.ForeignTableName]; ok {
			// skip already included preloads
			continue
		}
		sql := fmt.Sprintf(
			"LEFT JOIN %s ON %s.%s = %s.%s AND %s.deleted_at IS NULL",
			r.ForeignTableName, r.ForeignTableName, r.ForeignColumnName, r.TableName, r.ColumnName, r.ForeignTableName)
		(*d).Joins(sql)

		listCtx.groupBy = append(listCtx.groupBy, r.ForeignTableName+".id")
		listCtx.set[r.ForeignTableName] = true
	}
	if len(listCtx.joins) > 0 {
		// Add base relation
		listCtx.groupBy = append(listCtx.groupBy, (*d).GetTableName()+".id")
		(*d).Group(strings.Join(listCtx.groupBy, ","))
	}

	// Reset list of joins and group by's
	listCtx.joins = map[string]dao.TableRelation{}
}

func (s *sqlGenericService) loadList(listCtx *listContext, d *dao.GenericDao) *errors.ServiceError {
	args := listCtx.args
	ulog := *listCtx.ulog

	(*d).Count(listCtx.resourceList, &listCtx.pagingMeta.Total)

	// Set resourceList to be an empty slice with zero capacity. Real space will be allocated by g2.Find()
	if err := zeroSlice(listCtx.resourceList, 0); err != nil {
		return err
	}

	switch {
	case args.Size > MAX_LIST_SIZE:
		ulog.Warning("A query with a size greater than the maximum was requested.")
	case args.Size < 0:
		ulog.Warning("A query with an unbound size was requested.")
	case args.Size == 0:
		// This early return is not only performant, but also necessary.
		// gorm does not support Limit(0) any longer.
		ulog.Infof("A query with 0 size requested, returning early without collecting any resources from database")
		return nil
	}

	// NOTE: Limit no longer supports '0' size and will cause issues. There is an early return, do not remove it.
	//       https://github.com/go-gorm/gorm/blob/master/clause/limit.go#L18-L21
	if err := (*d).Fetch((args.Page-1)*int(args.Size), int(args.Size), listCtx.resourceList); err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			listCtx.pagingMeta.Size = 0
		} else {
			return errors.GeneralError("Unable to list resources: %s", err)
		}
	}
	listCtx.pagingMeta.Size = int64(reflect.ValueOf(listCtx.resourceList).Elem().Len())

	return nil
}

// Allocate a slice with size 'cap' of the type i
func zeroSlice(i interface{}, cap int64) *errors.ServiceError {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		return errors.GeneralError("A non-pointer to a list of resources: %v", v.Type())
	}
	// get the value that the pointer v points to.
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return errors.GeneralError("A non-slice list of resources")
	}
	v.Set(reflect.MakeSlice(v.Type(), 0, int(cap)))
	return nil
}

// walk the TSL tree looking for fields like, e.g., creator.username, and then:
// (1) look up the related table by its 1st part - creator
// (2) replace it by table name - creator.username -> accounts.username
func (s *sqlGenericService) treeWalkForRelatedTables(
	listCtx *listContext,
	tslTree tsl.Node,
	genericDao *dao.GenericDao,
) (tsl.Node, *errors.ServiceError) {
	resourceTable := (*genericDao).GetTableName()
	if listCtx.joins == nil {
		listCtx.joins = map[string]dao.TableRelation{}
	}
	walkFn := func(field string) (string, error) {
		fieldParts := strings.Split(field, ".")
		if len(fieldParts) > 1 && fieldParts[0] != resourceTable {
			nestedResource := fieldParts[0]
			_, exists := listCtx.joins[nestedResource]
			if !exists {
				// Populates relation if join exists
				if relation, ok := (*genericDao).GetTableRelation(nestedResource); ok {
					listCtx.joins[nestedResource] = relation
				} else if _, ok := listCtx.openapiToModel[field]; !ok {
					// If also not exposed as a nested resource consider this is an error
					return field, fmt.Errorf("%s is not a related resource of %s", strings.Join(fieldParts, "."), listCtx.resourceType)
				}
			}
			// replace by table name if coming from join
			if value, ok := listCtx.joins[nestedResource]; ok {
				fieldParts[0] = value.ForeignTableName
			}
			return strings.Join(fieldParts, "."), nil
		}
		return field, nil
	}

	tslTree, err := ident.Walk(tslTree, walkFn)
	if err != nil {
		return tslTree, errors.BadRequest(err.Error())
	}

	return tslTree, nil
}

// prepend table name to these "free" identifiers since they could cause "ambiguous" errors
func (s *sqlGenericService) treeWalkForAddingTableName(
	listCtx *listContext,
	tslTree tsl.Node,
	dao *dao.GenericDao,
) (tsl.Node, *errors.ServiceError) {
	resourceTable := (*dao).GetTableName()

	walkFn := func(field string) (string, error) {
		fieldParts := strings.Split(field, ".")
		if len(fieldParts) == 1 {
			if strings.Contains(field, "->") {
				return field, nil
			}
			return fmt.Sprintf("%s.%s", resourceTable, field), nil
		}
		return field, nil
	}

	tslTree, err := ident.Walk(tslTree, walkFn)
	if err != nil {
		return tslTree, errors.BadRequest(err.Error())
	}

	return tslTree, nil
}

func (s *sqlGenericService) treeWalkForSqlizer(
	listCtx *listContext,
	tslTree tsl.Node,
) (tsl.Node, squirrel.Sqlizer, *errors.ServiceError) {
	// Check field names in tree
	tslTree, serviceErr := db.FieldNameWalk(tslTree, listCtx.disallowedFields, listCtx.openapiToModel)
	if serviceErr != nil {
		return tslTree, nil, serviceErr
	}

	// Convert the search tree into SQL [Squirrel] filter
	sqlizer, err := sqlFilter.Walk(tslTree)
	if err != nil {
		return tslTree, nil, errors.BadRequest(err.Error())
	}

	return tslTree, sqlizer, nil
}
