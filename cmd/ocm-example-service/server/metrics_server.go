package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/handlers"
	"github.com/openshift-online/rh-trex/pkg/logger"
)

func NewMetricsServer() Server {
	mainRouter := mux.NewRouter()
	mainRouter.NotFoundHandler = http.HandlerFunc(api.SendNotFound)

	// metrics endpoint
	prometheusMetricsHandler := handlers.NewPrometheusMetricsHandler()
	mainRouter.Handle("/metrics", prometheusMetricsHandler.Handler())

	var mainHandler http.Handler = mainRouter

	s := &metricsServer{}
	s.httpServer = &http.Server{
		Addr:    env().Config.Metrics.BindAddress,
		Handler: mainHandler,
	}
	return s
}

type metricsServer struct {
	httpServer *http.Server
}

var _ Server = &metricsServer{}

func (s metricsServer) Listen() (listener net.Listener, err error) {
	return nil, nil
}

func (s metricsServer) Serve(listener net.Listener) {
}

func (s metricsServer) Start() {
	log := logger.NewOCMLogger(context.Background())
	var err error
	if env().Config.Metrics.EnableHTTPS {
		if env().Config.Server.HTTPSCertFile == "" || env().Config.Server.HTTPSKeyFile == "" {
			check(
				fmt.Errorf("Unspecified required --https-cert-file, --https-key-file"),
				"Can't start https server",
			)
		}

		// Serve with TLS
		log.Infof("Serving Metrics with TLS at %s", env().Config.Server.BindAddress)
		err = s.httpServer.ListenAndServeTLS(env().Config.Server.HTTPSCertFile, env().Config.Server.HTTPSKeyFile)
	} else {
		log.Infof("Serving Metrics without TLS at %s", env().Config.Metrics.BindAddress)
		err = s.httpServer.ListenAndServe()
	}
	check(err, "Metrics server terminated with errors")
	log.Infof("Metrics server terminated")
}

func (s metricsServer) Stop() error {
	return s.httpServer.Shutdown(context.Background())
}
