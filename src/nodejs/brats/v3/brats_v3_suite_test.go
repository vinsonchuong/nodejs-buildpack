package v3

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBratsV3(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brats V3 Suite")
}
