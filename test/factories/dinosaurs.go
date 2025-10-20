package factories

import (
	"context"
	"fmt"

	"github.com/openshift-online/rh-trex/cmd/trex/environments"
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/plugins/dinosaurs"
)

func (f *Factories) NewDinosaur(species string) (*api.Dinosaur, error) {
	dinoService := dinosaurs.Service(&environments.Environment().Services)

	dinosaur := &api.Dinosaur{
		Species: species,
	}

	dino, err := dinoService.Create(context.Background(), dinosaur)
	if err != nil {
		return nil, err
	}

	return dino, nil
}

func (f *Factories) NewDinosaurList(namePrefix string, count int) ([]*api.Dinosaur, error) {
	var dinosaurs []*api.Dinosaur
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("%s_%d", namePrefix, i)
		c, _ := f.NewDinosaur(name)
		dinosaurs = append(dinosaurs, c)
	}
	return dinosaurs, nil
}
