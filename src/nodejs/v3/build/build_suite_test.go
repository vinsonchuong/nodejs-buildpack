package build_test

import (
	. "github.com/onsi/ginkgo"
	"testing"
)

var T *testing.T

func TestSupply(t *testing.T) {
	RegisterFailHandler(Fail)
	T = t
	RunSpecs(t, "V3 Build Suite")
}
