package factories

import "github.com/segmentio/ksuid"

type Factories struct {
}

func (f *Factories) NewID() string {
	return ksuid.New().String()
}
