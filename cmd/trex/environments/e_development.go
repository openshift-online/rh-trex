package environments

import (
	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/openshift-online/rh-trex/pkg/db/db_session"
)

// devEnvImpl environment is intended for local use while developing features
type devEnvImpl struct {
	env *Env
}

var _ EnvironmentImpl = &devEnvImpl{}

func (e *devEnvImpl) OverrideDatabase(c *Database) error {
	c.SessionFactory = db_session.NewProdFactory(e.env.Config.Database)
	return nil
}

func (e *devEnvImpl) OverrideConfig(c *config.ApplicationConfig) error {
	c.Server.EnableJWT = false
	c.Server.EnableHTTPS = false
	return nil
}

func (e *devEnvImpl) OverrideServices(s *Services) error {
	return nil
}

func (e *devEnvImpl) OverrideHandlers(h *Handlers) error {
	return nil
}

func (e *devEnvImpl) OverrideClients(c *Clients) error {
	return nil
}

func (e *devEnvImpl) Flags() map[string]string {
	return map[string]string{
		"v":                      "10",
		"enable-authz":           "false",
		"ocm-debug":              "false",
		"enable-ocm-mock":        "true",
		"enable-https":           "false",
		"enable-metrics-https":   "false",
		"api-server-hostname":    "localhost",
		"api-server-bindaddress": "localhost:8000",
		"enable-sentry":          "false",
	}
}
