package servecmd

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/openshift-online/rh-trex/cmd/ocm-example-service/environments"
	"github.com/openshift-online/rh-trex/cmd/ocm-example-service/server"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the ocm-example-service",
		Long:  "Serve the ocm-example-service.",
		Run:   runServe,
	}
	err := environments.Environment().AddFlags(cmd.PersistentFlags())
	if err != nil {
		glog.Fatalf("Unable to add environment flags to serve command: %s", err.Error())
	}

	return cmd
}

func runServe(cmd *cobra.Command, args []string) {
	err := environments.Environment().Initialize()
	if err != nil {
		glog.Fatalf("Unable to initialize environment: %s", err.Error())
	}

	// Run the servers
	go func() {
		apiserver := server.NewAPIServer()
		apiserver.Start()
	}()

	go func() {
		metricsServer := server.NewMetricsServer()
		metricsServer.Start()
	}()

	go func() {
		healthcheckServer := server.NewHealthCheckServer()
		healthcheckServer.Start()
	}()

	go func() {
		controllersServer := server.NewControllersServer()
		controllersServer.Start()
	}()

	select {}
}
