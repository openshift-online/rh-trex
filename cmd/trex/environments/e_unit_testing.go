package environments

import (
	"os"

	"github.com/openshift-online/rh-trex/pkg/config"
	dbmocks "github.com/openshift-online/rh-trex/pkg/db/mocks"
)

var _ EnvironmentImpl = &unitTestingEnvImpl{}

// unitTestingEnvImpl is configuration for unit tests using mocked database
type unitTestingEnvImpl struct {
	env *Env
}

func (e *unitTestingEnvImpl) OverrideDatabase(c *Database) error {
	c.SessionFactory = dbmocks.NewMockSessionFactory()
	return nil
}

func (e *unitTestingEnvImpl) OverrideConfig(c *config.ApplicationConfig) error {
	// Support a one-off env to allow enabling db debug in testing
	if os.Getenv("DB_DEBUG") == "true" {
		c.Database.Debug = true
	}
	return nil
}

func (e *unitTestingEnvImpl) OverrideServices(s *Services) error {
	return nil
}

func (e *unitTestingEnvImpl) OverrideHandlers(h *Handlers) error {
	return nil
}

func (e *unitTestingEnvImpl) OverrideClients(c *Clients) error {
	return nil
}

func (e *unitTestingEnvImpl) Flags() map[string]string {
	return map[string]string{
		"v":                    "0",
		"logtostderr":          "true",
		"ocm-base-url":         "https://api.integration.openshift.com",
		"enable-https":         "false",
		"enable-metrics-https": "false",
		"enable-authz":         "true",
		"ocm-debug":            "false",
		"enable-ocm-mock":      "true",
		"enable-sentry":        "false",
	}
}
