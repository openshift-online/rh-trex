package util

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/rh-trex/pkg/util/test"
)

type SomeType struct{}

var _ = Describe("Strings util", func() {
	Describe("GetBaseType", func() {
		It("computes", func() {
			Expect(GetBaseType("asd")).To(Equal("string"))
			Expect(GetBaseType(SomeType{})).To(Equal("SomeType"))
			Expect(GetBaseType(&SomeType{})).To(Equal("SomeType"))
			Expect(GetBaseType(test.TestSomeType{})).To(Equal("TestSomeType"))
		})
	})

	Describe("GetType", func() {
		It("computes", func() {
			Expect(GetType("asd")).To(Equal("string"))
			Expect(GetType(SomeType{})).To(Equal("util.SomeType"))
			Expect(GetType(&SomeType{})).To(Equal("*util.SomeType"))
			Expect(GetType(test.TestSomeType{})).To(Equal("test.TestSomeType"))
		})
	})
})
