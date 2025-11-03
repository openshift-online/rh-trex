package environments

import (
	"os"

	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/openshift-online/rh-trex/pkg/db/db_session"
)

var _ EnvironmentImpl = &integrationTestingEnvImpl{}

// integrationTestingEnvImpl is configuration for integration tests using testcontainers
type integrationTestingEnvImpl struct {
	env *Env
}

func (e *integrationTestingEnvImpl) OverrideDatabase(c *Database) error {
	c.SessionFactory = db_session.NewTestcontainerFactory(e.env.Config.Database)
	return nil
}

func (e *integrationTestingEnvImpl) OverrideConfig(c *config.ApplicationConfig) error {
	// Support a one-off env to allow enabling db debug in testing
	if os.Getenv("DB_DEBUG") == "true" {
		c.Database.Debug = true
	}
	return nil
}

func (e *integrationTestingEnvImpl) OverrideServices(s *Services) error {
	return nil
}

func (e *integrationTestingEnvImpl) OverrideHandlers(h *Handlers) error {
	return nil
}

func (e *integrationTestingEnvImpl) OverrideClients(c *Clients) error {
	return nil
}

func (e *integrationTestingEnvImpl) Flags() map[string]string {
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
