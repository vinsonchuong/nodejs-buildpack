package integration_test

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CF NodeJS Buildpack", func() {
	var app *cutlass.App
	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	Context("deploying a Node app with a supported node version in .nvmrc and no version present in package.json engines", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "with_nvmrc_only"))
		})

		It("should supply the node version specified in .nvmrc", func() {
			PushAppAndConfirm(app)
			Expect(app.Stdout.String()).To(ContainSubstring("Downloading and installing npm 10.1.14"))
		})
	})

	Context("deploying a Node app with a supported node version in .nvmrc and an unsupported version in package.json engines", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "with_supported_nvmrc_unsupported_packagejson"))
		})

		It("should supply the node version specified in .nvmrc", func() {
			PushAppAndConfirm(app)
			Expect(app.Stdout.String()).To(ContainSubstring("Downloading and installing npm 10.1.14"))
		})
	})

	Context("deploying a Node app with the same supported node version in .nvmrc and package.json engines", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "with_supported_nvmrc_packagejson"))
		})

		It("should supply the node version specified in .nvmrc", func() {
			PushAppAndConfirm(app)
			Expect(app.Stdout.String()).To(ContainSubstring("Downloading and installing npm 10.1.14"))
		})
	})

	Context("deploying a Node app with different supported node versions in .nvmrc and package.json engines", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "with_different_supported_nvmrc_packagejson"))
		})

		It("should supply the node version that's most specific", func() {
			PushAppAndConfirm(app)
			Expect(app.Stdout.String()).To(ContainSubstring("Downloading and installing npm 10.1.14"))
		})
	})

	Context("deploying a Node app with an unsupported node version in .nvmrc and no version present in package.json engines", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "with_unsupported_nvmrc"))
		})

		It("fall back to default node version", func() {
			PushAppAndConfirm(app)
			Expect(app.Stdout.String()).To(ContainSubstring("Downloading and installing npm 10.1.14"))
		})
	})

	Context("deploying a Node app with an unsupported node version in .nvmrc and a supported version in package.json engines", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "with_unsupported_nvmrc_supported_packagejson"))
		})

		It("should supply the node version in package.json engines", func() {
			PushAppAndConfirm(app)
			Expect(app.Stdout.String()).To(ContainSubstring("Downloading and installing npm 10.1.14"))
		})
	})

})
