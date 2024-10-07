package services

import (
	"context"
	"net/url"
	"reflect"

	"github.com/openshift-online/rh-trex/pkg/auth"
	"github.com/openshift-online/rh-trex/pkg/client/ocm"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
	"go.uber.org/mock/gomock"

	"github.com/onsi/gomega/types"
	"github.com/yaacov/tree-search-language/pkg/tsl"

	azv1 "github.com/openshift-online/ocm-sdk-go/authorizations/v1"
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/openshift-online/rh-trex/pkg/db/db_session"
	"github.com/openshift-online/rh-trex/pkg/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type GenericTestDinosaur struct {
	api.Meta
	Species string
	// This is to illustrate resource review in action
	// It passes integration tests as it's mocked
	// does not work for local envs pointing to integration AMS via proxy
	OrganizationId string
}

var _ = Describe("populates search restriction", func() {
	var ctx context.Context
	var ctrl *gomock.Controller
	var genericService sqlGenericService
	var genericDao dao.GenericDao
	var authorizationMock *ocm.MockOCMAuthorization
	var ocmClientMock *ocm.Client
	username := "test-user"
	BeforeEach(func() {
		ctx = context.Background()
		ctx = auth.SetUsernameContext(ctx, username)
		ctrl = gomock.NewController(GinkgoT())
		dbConfig := config.NewDatabaseConfig()
		err := dbConfig.ReadFiles()
		Expect(err).ToNot(HaveOccurred())
		var dbFactory db.SessionFactory = db_session.NewTestFactory(dbConfig)
		defer dbFactory.Close()

		authorizationMock = ocm.NewMockOCMAuthorization(ctrl)
		ocmClientMock = &ocm.Client{
			Authorization: authorizationMock,
		}
		genericDao = dao.NewGenericDao(&dbFactory)
		genericService = sqlGenericService{
			genericDao: genericDao,
			ocmClient:  ocmClientMock,
		}
	})
	Context("Resource includes organization ID field", func() {
		When("Auth allows all orgs", func() {
			It("Allows all orgs", func() {
				args := NewListArguments(url.Values{})
				listCtx, model, serviceErr := newListContext(ctx, args, &[]GenericTestDinosaur{})
				resourceModel := reflect.TypeOf(&GenericTestDinosaur{}).Elem()
				Expect(model).To(Equal(reflect.New(resourceModel).Interface()))
				Expect(serviceErr).ToNot(HaveOccurred())
				response, err := azv1.NewResourceReview().
					AccountUsername(listCtx.username).
					Action(auth.GetAction).
					ResourceType("GenericTestDinosaur").
					OrganizationIDs("*").
					Build()
				Expect(err).ToNot(HaveOccurred())
				authorizationMock.EXPECT().
					ResourceReview(listCtx.ctx, listCtx.username, auth.GetAction, "GenericTestDinosaur").
					Return(response, nil)
				serviceErr = genericService.populateSearchRestriction(listCtx, model)
				Expect(serviceErr).ToNot(HaveOccurred())
				Expect(listCtx.args.Search).To(BeEmpty())
			})
		})
		When("Auth restricts orgs", func() {
			It("Allows only returned orgs", func() {
				args := NewListArguments(url.Values{})
				listCtx, model, serviceErr := newListContext(ctx, args, &[]GenericTestDinosaur{})
				resourceModel := reflect.TypeOf(&GenericTestDinosaur{}).Elem()
				Expect(model).To(Equal(reflect.New(resourceModel).Interface()))
				Expect(serviceErr).ToNot(HaveOccurred())
				response, err := azv1.NewResourceReview().
					AccountUsername(listCtx.username).
					Action(auth.GetAction).
					ResourceType("GenericTestDinosaur").
					OrganizationIDs("123", "124").
					Build()
				Expect(err).ToNot(HaveOccurred())
				authorizationMock.EXPECT().
					ResourceReview(listCtx.ctx, listCtx.username, auth.GetAction, "GenericTestDinosaur").
					Return(response, nil)
				serviceErr = genericService.populateSearchRestriction(listCtx, model)
				Expect(serviceErr).ToNot(HaveOccurred())
				Expect(listCtx.args.Search).ToNot(BeEmpty())
				Expect(listCtx.args.Search).To(Equal("organization_id in ('123','124')"))
			})
			It("Includes pre existing search", func() {
				args := NewListArguments(url.Values{})
				args.Search = "justification like '%test%'"
				listCtx, model, serviceErr := newListContext(ctx, args, &[]GenericTestDinosaur{})
				resourceModel := reflect.TypeOf(&GenericTestDinosaur{}).Elem()
				Expect(model).To(Equal(reflect.New(resourceModel).Interface()))
				Expect(serviceErr).ToNot(HaveOccurred())
				response, err := azv1.NewResourceReview().
					AccountUsername(listCtx.username).
					Action(auth.GetAction).
					ResourceType("GenericTestDinosaur").
					OrganizationIDs("123", "124").
					Build()
				Expect(err).ToNot(HaveOccurred())
				authorizationMock.EXPECT().
					ResourceReview(listCtx.ctx, listCtx.username, auth.GetAction, "GenericTestDinosaur").
					Return(response, nil)
				serviceErr = genericService.populateSearchRestriction(listCtx, model)
				Expect(serviceErr).ToNot(HaveOccurred())
				Expect(listCtx.args.Search).ToNot(BeEmpty())
				Expect(
					listCtx.args.Search,
				).To(Equal("justification like '%test%' and organization_id in ('123','124')"))
			})
		})
	})
})

var _ = Describe("Sql Translation", func() {
	var genericService sqlGenericService
	var genericDao dao.GenericDao
	BeforeEach(func() {
		dbConfig := config.NewDatabaseConfig()
		err := dbConfig.ReadFiles()
		Expect(err).ToNot(HaveOccurred())
		var dbFactory db.SessionFactory = db_session.NewTestFactory(dbConfig)
		defer dbFactory.Close()

		genericDao = dao.NewGenericDao(&dbFactory)
		genericService = sqlGenericService{genericDao: genericDao}
	})
	DescribeTable(
		"Errors",
		func(
			search string, errorMsg string) {
			listCtx, model, serviceErr := newListContext(
				context.Background(),
				&ListArguments{Search: search},
				&[]api.Dinosaur{},
			)
			Expect(serviceErr).ToNot(HaveOccurred())
			d := genericDao.GetInstanceDao(context.Background(), model)
			listCtx.disallowedFields = []string{"dinosaurs.id"}
			_, serviceErr = genericService.buildSearch(listCtx, &d)
			Expect(serviceErr).To(HaveOccurred())
			Expect(serviceErr.Code).To(Equal(errors.ErrorBadRequest))
			Expect(serviceErr.Error()).To(Equal(errorMsg))
		},
		Entry("Garbage", "garbage", "rh-trex-21: Failed to parse search query: garbage"),
		Entry("Disallowed field name", "id in ('123')", "rh-trex-21: dinosaurs.id is a disallowed field name"),
		Entry("Unknown field name", "bike = '123'", "rh-trex-21: dinosaurs.bike is not a valid field name"),
		Entry(
			"Unknown relation field",
			"status.bike = '123'",
			"rh-trex-21: status.bike is not a related resource of Dinosaur",
		),
	)

	DescribeTable(
		"Sql Parsing",
		func(
			search string, sqlReal string, valuesReal types.GomegaMatcher) {
			listCtx, _, serviceErr := newListContext(
				context.Background(),
				&ListArguments{Search: search},
				&[]api.Dinosaur{},
			)
			Expect(serviceErr).ToNot(HaveOccurred())
			tslTree, err := tsl.ParseTSL(search)
			Expect(err).ToNot(HaveOccurred())
			_, sqlizer, serviceErr := genericService.treeWalkForSqlizer(listCtx, tslTree)
			Expect(serviceErr).ToNot(HaveOccurred())
			sql, values, err := sqlizer.ToSql()
			Expect(err).ToNot(HaveOccurred())
			Expect(sql).To(Equal(sqlReal))
			Expect(values).To(valuesReal)
		},
		Entry(
			"Valid search",
			"dinosaurs.species like '%test%'",
			"dinosaurs.species LIKE ?",
			ConsistOf("%test%"),
		),
	)
})
