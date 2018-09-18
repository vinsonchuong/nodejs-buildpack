package build_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSupply(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V3 build Suite")
}
