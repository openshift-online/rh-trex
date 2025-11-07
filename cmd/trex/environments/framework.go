package environments

import (
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/pflag"

	"github.com/openshift-online/rh-trex/cmd/trex/environments/registry"
	"github.com/openshift-online/rh-trex/pkg/client/ocm"
	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

func init() {
	once.Do(func() {
		environment = &Env{}

		// Create the configuration
		environment.Config = config.NewApplicationConfig()
		environment.Name = GetEnvironmentStrFromEnv()

		environments = map[string]EnvironmentImpl{
			DevelopmentEnv:        &devEnvImpl{environment},
			UnitTestingEnv:        &unitTestingEnvImpl{environment},
			IntegrationTestingEnv: &integrationTestingEnvImpl{environment},
			ProductionEnv:         &productionEnvImpl{environment},
		}
	})
}

// EnvironmentImpl defines a set of behaviors for an OCM environment.
// Each environment provides a set of flags for basic set/override of the environment
// and configuration functions for each component type.
type EnvironmentImpl interface {
	Flags() map[string]string
	OverrideConfig(c *config.ApplicationConfig) error
	OverrideServices(s *Services) error
	OverrideDatabase(s *Database) error
	OverrideHandlers(c *Handlers) error
	OverrideClients(c *Clients) error
}

func GetEnvironmentStrFromEnv() string {
	envStr, specified := os.LookupEnv(EnvironmentStringKey)
	if !specified || envStr == "" {
		envStr = EnvironmentDefault
	}
	return envStr
}

func Environment() *Env {
	return environment
}

// AddFlags Adds environment flags, using the environment's config struct, to the flagset 'flags'
func (e *Env) AddFlags(flags *pflag.FlagSet) error {
	e.Config.AddFlags(flags)
	return setConfigDefaults(flags, environments[e.Name].Flags())
}

// Initialize loads the environment's resources
// This should be called after the e.Config has been set appropriately though AddFlags and pasing, done elsewhere
// The environment does NOT handle flag parsing
func (e *Env) Initialize() error {
	glog.Infof("Initializing %s environment", e.Name)

	envImpl, found := environments[e.Name]
	if !found {
		glog.Fatalf("Unknown runtime environment: %s", e.Name)
	}

	if err := envImpl.OverrideConfig(e.Config); err != nil {
		glog.Fatalf("Failed to configure ApplicationConfig: %s", err)
	}

	messages := environment.Config.ReadFiles()
	if len(messages) != 0 {
		glog.Fatalf("unable to read configuration files:\n%s", strings.Join(messages, "\n"))
	}

	// each env will set db explicitly because the DB impl has a `once` init section
	if err := envImpl.OverrideDatabase(&e.Database); err != nil {
		glog.Fatalf("Failed to configure Database: %s", err)
	}

	err := e.LoadClients()
	if err != nil {
		return err
	}
	if err := envImpl.OverrideClients(&e.Clients); err != nil {
		glog.Fatalf("Failed to configure Clients: %s", err)
	}

	e.LoadServices()
	if err := envImpl.OverrideServices(&e.Services); err != nil {
		glog.Fatalf("Failed to configure Services: %s", err)
	}

	seedErr := e.Seed()
	if seedErr != nil {
		return seedErr
	}

	if err := envImpl.OverrideHandlers(&e.Handlers); err != nil {
		glog.Fatalf("Failed to configure Handlers: %s", err)
	}

	return nil
}

func (e *Env) Seed() *errors.ServiceError {
	return nil
}

func (e *Env) LoadServices() {
	// Initialize the service registry map
	e.Services.serviceRegistry = make(map[string]interface{})

	// Auto-discovered services (no manual editing needed)
	registry.LoadDiscoveredServices(&e.Services, e)
}

func (e *Env) LoadClients() error {
	var err error

	ocmConfig := ocm.Config{
		BaseURL:      e.Config.OCM.BaseURL,
		ClientID:     e.Config.OCM.ClientID,
		ClientSecret: e.Config.OCM.ClientSecret,
		SelfToken:    e.Config.OCM.SelfToken,
		TokenURL:     e.Config.OCM.TokenURL,
		Debug:        e.Config.OCM.Debug,
	}

	// Create OCM Authz client
	if e.Config.OCM.EnableMock {
		glog.Infof("Using Mock OCM Authz Client")
		e.Clients.OCM, err = ocm.NewClientMock(ocmConfig)
	} else {
		e.Clients.OCM, err = ocm.NewClient(ocmConfig)
	}
	if err != nil {
		glog.Errorf("Unable to create OCM Authz client: %s", err.Error())
		return err
	}

	return nil
}



func (e *Env) Teardown() {
	if e.Database.SessionFactory != nil {
		if err := e.Database.SessionFactory.Close(); err != nil {
			glog.Errorf("Error closing database session factory: %s", err.Error())
		}
	}
	e.Clients.OCM.Close()
}

func setConfigDefaults(flags *pflag.FlagSet, defaults map[string]string) error {
	for name, value := range defaults {
		if err := flags.Set(name, value); err != nil {
			glog.Errorf("Error setting flag %s: %v", name, err)
			return err
		}
	}
	return nil
}
