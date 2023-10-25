package presenters

import (
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/api"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/api/openapi"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/util"
)

func ConvertDinosaur(dinosaur openapi.Dinosaur) *api.Dinosaur {
	return &api.Dinosaur{
		Meta: api.Meta{
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
