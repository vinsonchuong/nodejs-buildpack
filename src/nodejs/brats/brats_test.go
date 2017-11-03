package brats_test

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"nodejs/brats/helper"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/cloudfoundry/libbuildpack"
	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nodejs buildpack", func() {
	var app *cutlass.App
	AfterEach(func() { app = DestroyApp(app) })

	Context("Unbuilt buildpack (eg github)", func() {
		var bpName string
		BeforeEach(func() {
			bpName = GenBpName("unbuilt")
			cmd := exec.Command("git", "archive", "-o", filepath.Join("/tmp", bpName+".zip"), "HEAD")
			cmd.Dir = bpDir
			Expect(cmd.Run()).To(Succeed())
			Expect(cutlass.CreateOrUpdateBuildpack(bpName, filepath.Join("/tmp", bpName+".zip"))).To(Succeed())
			Expect(os.Remove(filepath.Join("/tmp", bpName+".zip"))).To(Succeed())

			app = cutlass.New(filepath.Join(bpDir, "fixtures", "brats"))
			app.Buildpacks = []string{bpName + "_buildpack"}
		})
		AfterEach(func() {
			Expect(cutlass.DeleteBuildpack(bpName)).To(Succeed())
		})

		It("runs", func() {
			PushApp(app)
			Expect(app.Stdout.String()).To(ContainSubstring("-----> Download go "))

			Expect(app.Stdout.String()).To(ContainSubstring("Installing node"))
			Expect(app.GetBody("/")).To(ContainSubstring("Hello World!"))
		})
	})

	Describe("deploying an app with an updated version of the same buildpack", func() {
		var bpName string
		BeforeEach(func() {
			bpName = GenBpName("changing")

			app = cutlass.New(filepath.Join(bpDir, "fixtures", "brats"))
			app.Buildpacks = []string{bpName + "_buildpack"}
		})
		AfterEach(func() {
			Expect(cutlass.DeleteBuildpack(bpName)).To(Succeed())
		})

		It("prints useful warning message to stdout", func() {
			Expect(cutlass.CreateOrUpdateBuildpack(bpName, buildpacks.UncachedFile)).To(Succeed())
			PushApp(app)
			Expect(app.Stdout.String()).ToNot(ContainSubstring("buildpack version changed from"))

			newFile, err := helper.ModifyBuildpack(buildpacks.UncachedFile, func(path string, r io.Reader) (io.Reader, error) {
				if path == "VERSION" {
					return strings.NewReader("NewVersion"), nil
				}
				return r, nil
			})
			Expect(err).ToNot(HaveOccurred())

			Expect(cutlass.CreateOrUpdateBuildpack(bpName, newFile)).To(Succeed())
			PushApp(app)
			Expect(app.Stdout.String()).To(MatchRegexp(`buildpack version changed from (\S+) to NewVersion`))
		})
	})

	Describe("For all supported Node versions", func() {
		bpDir, err := cutlass.FindRoot()
		if err != nil {
			panic(err)
		}
		manifest, err := libbuildpack.NewManifest(bpDir, nil, time.Now())
		nodeVersions := manifest.AllDependencyVersions("node")
		var appDir string
		AfterEach(func() { os.RemoveAll(appDir) })

		for _, v := range nodeVersions {
			nodeVersion := v
			It("with node "+nodeVersion, func() {
				appDir = CopyBrats(nodeVersion)
				app = cutlass.New(appDir)
				app.Buildpacks = []string{buildpacks.Cached}
				PushApp(app)

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
		}
	})

	Describe("staging with node buildpack that sets EOL on dependency", func() {
		var (
			eolDate       string
			buildpackFile string
			bpName        string
		)
		JustBeforeEach(func() {
			eolDate = time.Now().AddDate(0, 0, 10).Format("2006-01-02")
			file, err := helper.ModifyBuildpackManifest(buildpackFile, func(m *helper.Manifest) {
				for _, eol := range m.DependencyDeprecationDates {
					if eol.Name == "node" {
						eol.Date = eolDate
					}
				}
			})
			Expect(err).ToNot(HaveOccurred())
			bpName = GenBpName("eol")
			Expect(cutlass.CreateOrUpdateBuildpack(bpName, file)).To(Succeed())
			os.Remove(file)

			app = cutlass.New(filepath.Join(bpDir, "fixtures", "brats"))
			app.Buildpacks = []string{bpName + "_buildpack"}
			PushApp(app)
		})
		AfterEach(func() {
			Expect(cutlass.DeleteBuildpack(bpName)).To(Succeed())
		})

		Context("using an uncached buildpack", func() {
			BeforeEach(func() {
				buildpackFile = buildpacks.UncachedFile
			})
			It("warns about end of life", func() {
				Expect(app.Stdout.String()).To(MatchRegexp(`WARNING.*node \S+ will no longer be available in new buildpacks released after`))
			})
		})

		Context("using a cached buildpack", func() {
			BeforeEach(func() {
				buildpackFile = buildpacks.CachedFile
			})
			It("warns about end of life", func() {
				Expect(app.Stdout.String()).To(MatchRegexp(`WARNING.*node \S+ will no longer be available in new buildpacks released after`))
			})
		})
	})

	Describe("staging with a version of node that is not the latest patch release in the manifest", func() {
		var appDir string
		BeforeEach(func() {
			manifest, err := libbuildpack.NewManifest(bpDir, nil, time.Now())
			Expect(err).ToNot(HaveOccurred())
			raw := manifest.AllDependencyVersions("node")
			vs := make([]*semver.Version, len(raw))
			for i, r := range raw {
				vs[i], err = semver.NewVersion(r)
				Expect(err).ToNot(HaveOccurred())
			}
			sort.Sort(semver.Collection(vs))
			version := vs[0].Original()

			appDir = CopyBrats(version)
			app = cutlass.New(appDir)
			app.Buildpacks = []string{buildpacks.Cached}
			PushApp(app)
		})
		AfterEach(func() { os.RemoveAll(appDir) })

		It("logs a warning that tells the user to upgrade the dependency", func() {
			Expect(app.Stdout.String()).To(MatchRegexp("WARNING.*A newer version of node is available in this buildpack"))
		})
	})

	Describe("staging with custom buildpack that uses credentials in manifest dependency uris", func() {
		var (
			buildpackFile string
			bpName        string
		)
		JustBeforeEach(func() {
			file, err := helper.ModifyBuildpackManifest(buildpackFile, func(m *helper.Manifest) {
				for _, d := range m.Dependencies {
					uri, err := url.Parse(d.URI)
					uri.User = url.UserPassword("login", "password")
					Expect(err).ToNot(HaveOccurred())
					d.URI = uri.String()
				}
			})
			Expect(err).ToNot(HaveOccurred())
			bpName = GenBpName("eol")
			Expect(cutlass.CreateOrUpdateBuildpack(bpName, file)).To(Succeed())
			os.Remove(file)

			app = cutlass.New(filepath.Join(bpDir, "fixtures", "brats"))
			app.Buildpacks = []string{bpName + "_buildpack"}
			PushApp(app)
		})
		AfterEach(func() {
			Expect(cutlass.DeleteBuildpack(bpName)).To(Succeed())
		})
		Context("using an uncached buildpack", func() {
			BeforeEach(func() {
				buildpackFile = buildpacks.UncachedFile
			})
			It("does not include credentials in logged dependency uris", func() {
				Expect(app.Stdout.String()).To(MatchRegexp(`node\-[\d\.]+\-linux\-x64\-[\da-f]+\.tgz`))
				Expect(app.Stdout.String()).ToNot(ContainSubstring("login"))
				Expect(app.Stdout.String()).ToNot(ContainSubstring("password"))
			})
		})
		Context("using a cached buildpack", func() {
			BeforeEach(func() {
				buildpackFile = buildpacks.UncachedFile
			})
			It("does not include credentials in logged dependency file paths", func() {
				Expect(app.Stdout.String()).To(MatchRegexp(`node\-[\d\.]+\-linux\-x64\-[\da-f]+\.tgz`))
				Expect(app.Stdout.String()).ToNot(ContainSubstring("login"))
				Expect(app.Stdout.String()).ToNot(ContainSubstring("password"))
			})
		})
	})

	Describe("deploying an app that has an executable .profile script", func() {
		BeforeEach(func() {
			manifest, err := libbuildpack.NewManifest(bpDir, nil, time.Now())
			dep, err := manifest.DefaultVersion("node")
			Expect(err).ToNot(HaveOccurred())

			appDir := CopyBrats(dep.Version)
			AddDotProfileScriptToApp(appDir)
			app = cutlass.New(appDir)
			app.Buildpacks = []string{buildpacks.Cached}
			PushApp(app)
		})
		AfterEach(func() { os.RemoveAll(app.Path) })

		It("executes the .profile script", func() {
			Expect(app.Stdout.String()).To(ContainSubstring("PROFILE_SCRIPT_IS_PRESENT_AND_RAN"))
		})
		It("does not let me view the .profile script", func() {
			_, headers, err := app.Get("/.profile", map[string]string{})
			Expect(err).ToNot(HaveOccurred())
			Expect(headers).To(HaveKeyWithValue("StatusCode", []string{"404"}))
		})
	})

	Describe("deploying an app that has sensitive environment variables", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "brats"))
			app.Buildpacks = []string{buildpacks.Cached}
			app.SetEnv("MY_SPECIAL_VAR", "SUPER SENSITIVE DATA")
			PushApp(app)
		})

		It("will not write credentials to the app droplet", func() {
			Expect(app.DownloadDroplet(filepath.Join(app.Path, "droplet.tgz"))).To(Succeed())
			file, err := os.Open(filepath.Join(app.Path, "droplet.tgz"))
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()
			gz, err := gzip.NewReader(file)
			Expect(err).ToNot(HaveOccurred())
			defer gz.Close()
			tr := tar.NewReader(gz)

			for {
				hdr, err := tr.Next()
				if err == io.EOF {
					break
				}
				b, err := ioutil.ReadAll(tr)
				for _, content := range []string{"MY_SPECIAL_VAR", "SUPER SENSITIVE DATA"} {
					if strings.Contains(string(b), content) {
						Fail(fmt.Sprintf("Found sensitive string %s in %s", content, hdr.Name))
					}
				}
			}
		})
	})
})
