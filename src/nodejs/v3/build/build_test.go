package build_test

import (
	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack/test"
	. "github.com/onsi/ginkgo"
	"nodejs/v3/build"
)

var _ = Describe("MakeNodeLayers", func() {
	It("should fetch the node dependency and install it into the cache and launch layers", func() {
		buildFactory := test.NewBuildFactory(T)
		buildFactory.AddBuildPlan(T, "node", libbuildpack.BuildPlanDependency{})
		buildFactory.AddDependency(T, "node", "foo")

		build.MakeNodeLayers(buildFactory.Build)
	})
})
