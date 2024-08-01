package dao

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("applyBaseMapping", func() {
	It("generates base mapping", func() {
		result := map[string]string{}
		applyBaseMapping(result, []string{"id", "created_at", "column1", "nested.field"}, "test_table")
		for k, v := range result {
			if strings.HasPrefix(k, "test_table") {
				Expect(k).To(Equal(v))
				continue
			}
			// nested fields from table
			i := strings.Index(k, ".")
			Expect(k[i+1:]).To(Equal(v))
		}
	})
})

var _ = Describe("applyRelationMapping", func() {
	It("generates relation mapping", func() {
		result := map[string]string{}
		applyBaseMapping(result, []string{"id", "created_at", "column1", "nested.field"}, "base_table")
		applyRelationMapping(result, []relationMapping{
			func() TableMappingRelation {
				result := map[string]string{}
				applyBaseMapping(result, []string{"id", "created_at", "column1", "nested.field"}, "relation_table")
				return TableMappingRelation{
					relationTableName: "relation_table",
					Mapping:           result,
				}
			},
		})
		for k, v := range result {
			if strings.HasPrefix(k, "base_table") {
				Expect(k).To(Equal(v))
				continue
			}
			if strings.HasPrefix(k, "relation_table") {
				if c := strings.Count(k, "."); c > 1 {
					i := strings.Index(k, ".")
					i = strings.Index(k[i+1:], ".") + i
					Expect(k[i+2:]).To(Equal(v))
					continue
				}
				Expect(k).To(Equal(v))
				continue
			}

			// nested fields from base table
			i := strings.Index(k, ".")
			Expect(k[i+1:]).To(Equal(v))
			fmt.Println(k, v)
		}
	})
})
