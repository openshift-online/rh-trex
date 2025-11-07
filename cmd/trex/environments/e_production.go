package environments

import (
	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/openshift-online/rh-trex/pkg/db/db_session"
)

var _ EnvironmentImpl = &productionEnvImpl{}

// productionEnvImpl is any deployed instance of the service through app-interface
type productionEnvImpl struct {
	env *Env
}

func (e *productionEnvImpl) OverrideDatabase(c *Database) error {
	c.SessionFactory = db_session.NewProdFactory(e.env.Config.Database)
	return nil
}

func (e *productionEnvImpl) OverrideConfig(c *config.ApplicationConfig) error {
	return nil
}

func (e *productionEnvImpl) OverrideServices(s *Services) error {
	return nil
}

func (e *productionEnvImpl) OverrideHandlers(h *Handlers) error {
	return nil
}

func (e *productionEnvImpl) OverrideClients(c *Clients) error {
	return nil
}

func (e *productionEnvImpl) Flags() map[string]string {
	return map[string]string{
		"v":               "1",
		"ocm-debug":       "false",
		"enable-ocm-mock": "false",
	}
}
