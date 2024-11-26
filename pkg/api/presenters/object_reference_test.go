package presenters

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/rh-trex/pkg/api"
)

var _ = Describe("Object Reference Presenter", func() {
	It("Populates Kind", func() {
		object := api.Dinosaur{
			Meta: api.Meta{
				ID: "123",
			},
		}
		presented := PresentReference("123", object)
		Expect(*presented.Id).To(Equal(object.ID))
		Expect(*presented.Kind).To(Equal("Dinosaur"))
		Expect(*presented.Href).To(Equal("/api/rh-trex/v1/dinosaurs/123"))
	})
})
