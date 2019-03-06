package v2b_integration_test

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pushing an app a second time", func() {
	var app *cutlass.App
	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	BeforeEach(func() {
		//if cutlass.Cached {
		//	Skip("running uncached tests")
		//}

		app = cutlass.New(filepath.Join("testdata", "simple_app"))
		app.Buildpacks = []string{"nodejs_buildpack"}
	})

	//Regexp := `.*/node\-[\d\.]+\-linux\-x64\-(cflinuxfs.*-)?[\da-f]+\.tgz`
	//DownloadRegexp := "Download from " + Regexp
	//CopyRegexp := "Copy " + Regexp

	It("uses the cache for manifest dependencies", func() {
		PushAppAndConfirm(app)

		Expect(app.Stdout.String()).To(ContainSubstring("Installing node_modules"))
		Expect(app.Stdout.String()).NotTo(ContainSubstring("Reusing existing npm-cache"))

		app.Stdout.Reset()
		PushAppAndConfirm(app)

		Expect(app.Stdout.String()).To(ContainSubstring("Installing node_modules"))
		Expect(app.Stdout.String()).To(ContainSubstring("Reusing existing npm-cache"))

	})
})
