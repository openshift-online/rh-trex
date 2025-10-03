package factories

import (
	coreFactories "github.com/openshift-online/rh-trex-core/test/factories"
)

// Factories embeds the core factory to inherit common functionality
type Factories struct {
	coreFactories.Factories
}

// NewID creates a new unique ID using KSUID - now provided by core library
func (f *Factories) NewID() string {
	return f.Factories.NewID()
}
