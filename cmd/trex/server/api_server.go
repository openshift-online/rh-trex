package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	_ "github.com/auth0/go-jwt-middleware"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/ghodss/yaml"
	_ "github.com/golang-jwt/jwt/v4"
	"github.com/golang/glog"
	gorillahandlers "github.com/gorilla/handlers"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift-online/ocm-sdk-go/authentication"

	"github.com/openshift-online/rh-trex/cmd/trex/environments"
	"github.com/openshift-online/rh-trex/data/generated/openapi"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

type apiServer struct {
	httpServer *http.Server
}

var _ Server = &apiServer{}

func env() *environments.Env {
	return environments.Environment()
}

func NewAPIServer() Server {
	s := &apiServer{}

	mainRouter := s.routes()

	// Sentryhttp middleware performs two operations:
	// 1) Attaches an instance of *sentry.Hub to the request’s context. Accessit by using the sentry.GetHubFromContext() method on the request
	//   NOTE this is the only way middleware, handlers, and services should be reporting to sentry, through the hub
	// 2) Reports panics to the configured sentry service
	if env().Config.Sentry.Enabled {
		sentryhttpOptions := sentryhttp.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         env().Config.Sentry.Timeout,
		}
		sentryMW := sentryhttp.New(sentryhttpOptions)
		mainRouter.Use(sentryMW.Handle)
	}

	// referring to the router as type http.Handler allows us to add middleware via more handlers
	var mainHandler http.Handler = mainRouter

	if env().Config.Server.EnableJWT {
		// Create the logger for the authentication handler:
		authnLogger, err := sdk.NewGlogLoggerBuilder().
			InfoV(glog.Level(1)).
			DebugV(glog.Level(5)).
			Build()
		check(err, "Unable to create authentication logger")

		// Create the handler that verifies that tokens are valid:
		mainHandler, err = authentication.NewHandler().
			Logger(authnLogger).
			KeysFile(env().Config.Server.JwkCertFile).
			KeysURL(env().Config.Server.JwkCertURL).
			ACLFile(env().Config.Server.ACLFile).
			Public("^/api/rh-trex/?$").
			Public("^/api/rh-trex/v1/?$").
			Public("^/api/rh-trex/v1/openapi/?$").
			Public("^/api/rh-trex/v1/errors(/.*)?$").
			Next(mainHandler).
			Build()
		check(err, "Unable to create authentication handler")
	}

	// TODO: remove all cloud.redhat.com once migration to console.redhat.com is complete
	// refer to: https://issues.redhat.com/browse/RHCLOUD-14695
	mainHandler = gorillahandlers.CORS(
		gorillahandlers.AllowedOrigins([]string{
			// OCM UI local development URLs
			"https://qa.foo.redhat.com:1337",
			"https://prod.foo.redhat.com:1337",
			"https://ci.foo.redhat.com:1337",
			"https://cloud.redhat.com",   // TODO: remove
			"https://console.redhat.com", // Production / candidate
			// Staging and test environments
			"https://qaprodauth.cloud.redhat.com", // TODO: remove
			"https://qa.cloud.redhat.com",         // TODO: remove
			"https://ci.cloud.redhat.com",         // TODO: remove
			"https://qaprodauth.console.redhat.com",
			"https://qa.console.redhat.com",
			"https://ci.console.redhat.com",
			"https://console.stage.redhat.com",
			// API docs UI
			"https://api.stage.openshift.com",
			"https://api.openshift.com",
			// Customer portal
			"https://access.qa.redhat.com",
			"https://access.stage.redhat.com",
			"https://access.redhat.com",
		}),
		gorillahandlers.AllowedMethods([]string{
			http.MethodDelete,
			http.MethodGet,
			http.MethodPatch,
			http.MethodPost,
		}),
		gorillahandlers.AllowedHeaders([]string{
			"Authorization",
			"Content-Type",
		}),
		gorillahandlers.MaxAge(int((10 * time.Minute).Seconds())),
	)(mainHandler)

	mainHandler = removeTrailingSlash(mainHandler)

	s.httpServer = &http.Server{
		Addr:    env().Config.Server.BindAddress,
		Handler: mainHandler,
	}

	return s
}

// Serve start the blocking call to Serve.
// Useful for breaking up ListenAndServer (Start) when you require the server to be listening before continuing
func (s apiServer) Serve(listener net.Listener) {
	var err error
	if env().Config.Server.EnableHTTPS {
		// Check https cert and key path path
		if env().Config.Server.HTTPSCertFile == "" || env().Config.Server.HTTPSKeyFile == "" {
			check(
				fmt.Errorf("Unspecified required --https-cert-file, --https-key-file"),
				"Can't start https server",
			)
		}

		// Serve with TLS
		glog.Infof("Serving with TLS at %s", env().Config.Server.BindAddress)
		err = s.httpServer.ServeTLS(listener, env().Config.Server.HTTPSCertFile, env().Config.Server.HTTPSKeyFile)
	} else {
		glog.Infof("Serving without TLS at %s", env().Config.Server.BindAddress)
		err = s.httpServer.Serve(listener)
	}

	// Web server terminated.
	check(err, "Web server terminated with errors")
	glog.Info("Web server terminated")
}

// Listen only start the listener, not the server.
// Useful for breaking up ListenAndServer (Start) when you require the server to be listening before continuing
func (s apiServer) Listen() (listener net.Listener, err error) {
	return net.Listen("tcp", env().Config.Server.BindAddress)
}

// Start listening on the configured port and start the server. This is a convenience wrapper for Listen() and Serve(listener Listener)
func (s apiServer) Start() {
	listener, err := s.Listen()
	if err != nil {
		glog.Fatalf("Unable to start API server: %s", err)
	}
	s.Serve(listener)

	// after the server exits but before the application terminates
	// we need to explicitly close Go's sql connection pool.
	// this needs to be called *exactly* once during an app's lifetime.
	env().Database.SessionFactory.Close()
}

func (s apiServer) Stop() error {
	return s.httpServer.Shutdown(context.Background())
}

func (s *apiServer) loadOpenAPISpec(asset string) (data []byte, err error) {
	data, err = openapi.Asset(asset)
	if err != nil {
		err = errors.GeneralError(
			"can't load OpenAPI specification from asset '%s'",
			asset,
		)
		return
	}
	data, err = yaml.YAMLToJSON(data)
	if err != nil {
		err = errors.GeneralError(
			"can't convert OpenAPI specification loaded from asset '%s' from YAML to JSON",
			asset,
		)
		return
	}
	return
}
