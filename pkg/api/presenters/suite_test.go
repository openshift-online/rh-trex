package presenters

import (
	"testing"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

func TestAccessProtection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Presenters Suite")
}
