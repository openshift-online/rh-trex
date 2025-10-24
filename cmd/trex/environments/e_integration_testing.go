package environments

import (
	"os"

	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/openshift-online/rh-trex/pkg/db/db_session"
)

var _ EnvironmentImpl = &testingEnvImpl{}

// testingEnvImpl is configuration for local integration tests
type testingEnvImpl struct {
	env *Env
}

func (e *testingEnvImpl) OverrideDatabase(c *Database) error {
	c.SessionFactory = db_session.NewTestFactory(e.env.Config.Database)
	return nil
}

func (e *testingEnvImpl) OverrideConfig(c *config.ApplicationConfig) error {
	// Support a one-off env to allow enabling db debug in testing
	if os.Getenv("DB_DEBUG") == "true" {
		c.Database.Debug = true
	}
	return nil
}

func (e *testingEnvImpl) OverrideServices(s *Services) error {
	return nil
}

func (e *testingEnvImpl) OverrideHandlers(h *Handlers) error {
	return nil
}

func (e *testingEnvImpl) OverrideClients(c *Clients) error {
	return nil
}

func (e *testingEnvImpl) Flags() map[string]string {
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
