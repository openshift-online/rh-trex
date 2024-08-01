package util

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Strings util", func() {
	Describe("ToSnakeCase", func() {
		It("transforms", func() {
			Expect(ToSnakeCase("asd")).To(Equal("asd"))
			Expect(ToSnakeCase("AsdAsd")).To(Equal("asd_asd"))
			Expect(ToSnakeCase("asdAsd")).To(Equal("asd_asd"))
			Expect(ToSnakeCase("Asd Asd")).To(Equal("asd_asd"))
		})
	})
})
