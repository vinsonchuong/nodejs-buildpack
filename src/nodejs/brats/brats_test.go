package brats_test

import (
	"nodejs/brats/helper"

	"github.com/cloudfoundry/libbuildpack/cutlass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nodejs buildpack", func() {
	helper.UnbuiltBuildpack("node", CopyBrats)
	helper.DeployingAnAppWithAnUpdatedVersionOfTheSameBuildpack(CopyBrats)
	helper.StagingWithBuildpackThatSetsEOL("node", CopyBrats)
	helper.StagingWithADepThatIsNotTheLatest("node", CopyBrats)
	helper.StagingWithCustomBuildpackWithCredentialsInDependencies(`node\-[\d\.]+\-linux\-x64\-[\da-f]+\.tgz`, CopyBrats)
	helper.DeployAppWithExecutableProfileScript("node", CopyBrats)
	helper.DeployAnAppWithSensitiveEnvironmentVariables(CopyBrats)
	helper.ForAllSupportedVersions("node", CopyBrats, func(nodeVersion string, app *cutlass.App) {
		By("runs a simple webserver", func() {
			Expect(app.GetBody("/")).To(ContainSubstring("Hello World!"))
		})
		By("supports bcrypt", func() {
			Expect(app.GetBody("/bcrypt")).To(ContainSubstring("Hello Bcrypt!"))
		})
		By("supports bson-ext", func() {
			Expect(app.GetBody("/bson-ext")).To(ContainSubstring("Hello Bson-ext!"))
		})
		By("installs the correct version", func() {
			Expect(app.Stdout.String()).To(ContainSubstring("Installing node " + nodeVersion))
		})
	})
})
