package server

import (
	"fmt"
	"net/http"

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/openshift-online/rh-trex/cmd/trex/server/logging"
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/auth"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/pkg/handlers"
	"github.com/openshift-online/rh-trex/pkg/logger"
)

type ServicesInterface interface {
	GetService(name string) interface{}
}

type RouteRegistrationFunc func(apiV1Router *mux.Router, services ServicesInterface, authMiddleware auth.JWTMiddleware, authzMiddleware auth.AuthorizationMiddleware)

var routeRegistry = make(map[string]RouteRegistrationFunc)

func RegisterRoutes(name string, registrationFunc RouteRegistrationFunc) {
	routeRegistry[name] = registrationFunc
}

func LoadDiscoveredRoutes(apiV1Router *mux.Router, services ServicesInterface, authMiddleware auth.JWTMiddleware, authzMiddleware auth.AuthorizationMiddleware) {
	for name, registrationFunc := range routeRegistry {
		registrationFunc(apiV1Router, services, authMiddleware, authzMiddleware)
		_ = name // prevent unused variable warning
	}
}

func (s *apiServer) routes() *mux.Router {
	services := &env().Services

	metadataHandler := handlers.NewMetadataHandler()

	var authMiddleware auth.JWTMiddleware
	authMiddleware = &auth.MiddlewareMock{}
	if env().Config.Server.EnableJWT {
		var err error
		authMiddleware, err = auth.NewAuthMiddleware()
		check(err, "Unable to create auth middleware")
	}
	if authMiddleware == nil {
		check(fmt.Errorf("auth middleware is nil"), "Unable to create auth middleware: missing middleware")
	}

	authzMiddleware := auth.NewAuthzMiddlewareMock()
	if env().Config.Server.EnableAuthz {
		// TODO: authzMiddleware, err = auth.NewAuthzMiddleware()
		// check(err, "Unable to create authz middleware")
	}

	// mainRouter is top level "/"
	mainRouter := mux.NewRouter()
	mainRouter.NotFoundHandler = http.HandlerFunc(api.SendNotFound)

	// Operation ID middleware sets a relatively unique operation ID in the context of each request for debugging purposes
	mainRouter.Use(logger.OperationIDMiddleware)

	// Request logging middleware logs pertinent information about the request and response
	mainRouter.Use(logging.RequestLoggingMiddleware)

	//  /api/rh-trex
	apiRouter := mainRouter.PathPrefix("/api/rh-trex").Subrouter()
	apiRouter.HandleFunc("", metadataHandler.Get).Methods(http.MethodGet)

	//  /api/rh-trex/v1
	apiV1Router := apiRouter.PathPrefix("/v1").Subrouter()

	//  /api/rh-trex/v1/openapi
	openapiUIHandler, err := handlers.NewOpenAPIUIHandler()
	check(err, "Unable to create OpenAPI UI handler")
	apiV1Router.HandleFunc("/openapi", openapiUIHandler.Get).Methods(http.MethodGet)

	openapiHandler, err := handlers.NewOpenAPIHandler()
	check(err, "Unable to create OpenAPI handler")
	apiV1Router.HandleFunc("/openapi.json", openapiHandler.Get).Methods(http.MethodGet)
	registerApiMiddleware(apiV1Router)

	// Auto-discovered routes (no manual editing needed)
	LoadDiscoveredRoutes(apiV1Router, services, authMiddleware, authzMiddleware)

	return mainRouter
}

func registerApiMiddleware(router *mux.Router) {
	router.Use(MetricsMiddleware)

	router.Use(
		func(next http.Handler) http.Handler {
			return db.TransactionMiddleware(next, env().Database.SessionFactory)
		},
	)

	router.Use(gorillahandlers.CompressHandler)
}
