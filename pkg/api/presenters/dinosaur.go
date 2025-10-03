package presenters

import (
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/api/openapi"
	coreapi "github.com/openshift-online/rh-trex-core/api"
	"github.com/openshift-online/rh-trex/pkg/util"
)

func ConvertDinosaur(dinosaur openapi.Dinosaur) *api.Dinosaur {
	return &api.Dinosaur{
		Meta: coreapi.Meta{
			ID: util.NilToEmptyString(dinosaur.Id),
		},
		Species: util.NilToEmptyString(dinosaur.Species),
	}
}

func PresentDinosaur(dinosaur *api.Dinosaur) openapi.Dinosaur {
	reference := PresentReference(dinosaur.ID, dinosaur)
	return openapi.Dinosaur{
		Id:        reference.Id,
		Kind:      reference.Kind,
		Href:      reference.Href,
		Species:   openapi.PtrString(dinosaur.Species),
		CreatedAt: openapi.PtrTime(dinosaur.CreatedAt),
		UpdatedAt: openapi.PtrTime(dinosaur.UpdatedAt),
	}
}
