package brats_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cloudfoundry/libbuildpack/cutlass"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var bpDir string
var buildpacks struct {
	BpVersion    string
	BpLanguage   string
	Cached       string
	CachedFile   string
	Uncached     string
	UncachedFile string
}

func init() {
	flag.StringVar(&cutlass.DefaultMemory, "memory", "128M", "default memory for pushed apps")
	flag.StringVar(&cutlass.DefaultDisk, "disk", "256M", "default disk for pushed apps")
	flag.Parse()
}

var _ = SynchronizedBeforeSuite(func() []byte {
	// Run once
	buildpacks.BpVersion = cutlass.RandStringRunes(6)
	buildpacks.Cached = "brats_nodejs_cached_" + buildpacks.BpVersion
	buildpacks.Uncached = "brats_nodejs_uncached_" + buildpacks.BpVersion

	var err error
	bpDir, err = cutlass.FindRoot()
	Expect(err).NotTo(HaveOccurred())

	file, err := ioutil.ReadFile(filepath.Join(bpDir, "manifest.yml"))
	Expect(err).ToNot(HaveOccurred())
	obj := make(map[string]interface{})
	Expect(yaml.Unmarshal(file, &obj)).To(Succeed())
	var ok bool
	buildpacks.BpLanguage, ok = obj["language"].(string)
	Expect(ok).To(BeTrue())

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Fprintln(os.Stderr, "Start build cached buildpack")
		cachedBuildpack, err := cutlass.PackageUniquelyVersionedBuildpackExtra(buildpacks.Cached, buildpacks.BpVersion, true)
		Expect(err).NotTo(HaveOccurred())
		buildpacks.CachedFile = cachedBuildpack.File
		fmt.Fprintln(os.Stderr, "Finish cached buildpack")
	}()
	go func() {
		defer wg.Done()
		fmt.Fprintln(os.Stderr, "Start build uncached buildpack")
		uncachedBuildpack, err := cutlass.PackageUniquelyVersionedBuildpackExtra(buildpacks.Uncached, buildpacks.BpVersion, false)
		Expect(err).NotTo(HaveOccurred())
		buildpacks.UncachedFile = uncachedBuildpack.File
		fmt.Fprintln(os.Stderr, "Finish uncached buildpack")
	}()
	wg.Wait()

	buildpacks.Cached = buildpacks.Cached + "_buildpack"
	buildpacks.Uncached = buildpacks.Uncached + "_buildpack"

	// Marshall for run all nodes
	data, err := json.Marshal(buildpacks)
	Expect(err).NotTo(HaveOccurred())

	return data
}, func(data []byte) {
	// Run on all nodes
	err := json.Unmarshal(data, &buildpacks)
	Expect(err).NotTo(HaveOccurred())

	bpDir, err = cutlass.FindRoot()
	Expect(err).NotTo(HaveOccurred())

	cutlass.SeedRandom()
	cutlass.DefaultStdoutStderr = GinkgoWriter
})

var _ = SynchronizedAfterSuite(func() {
	// Run on all nodes
}, func() {
	// Run once
	Expect(cutlass.DeleteOrphanedRoutes()).To(Succeed())
	Expect(cutlass.DeleteBuildpack(strings.Replace(buildpacks.Cached, "_buildpack", "", 1))).To(Succeed())
	Expect(cutlass.DeleteBuildpack(strings.Replace(buildpacks.Uncached, "_buildpack", "", 1))).To(Succeed())
	Expect(os.Remove(buildpacks.CachedFile)).To(Succeed())
})

func TestBrats(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brats Suite")
}

func PushApp(app *cutlass.App) {
	Expect(app.Push()).To(Succeed())
	Eventually(app.InstanceStates, 20*time.Second).Should(Equal([]string{"RUNNING"}))
}

func DestroyApp(app *cutlass.App) *cutlass.App {
	if app != nil {
		app.Destroy()
	}
	return nil
}

func GenBpName(name string) string {
	return "brats_" + buildpacks.BpLanguage + "_" + name + "_" + cutlass.RandStringRunes(6)
}

func CopyBrats(nodejsVersion string) string {
	dir, err := cutlass.CopyFixture(filepath.Join(bpDir, "fixtures", "brats"))
	Expect(err).ToNot(HaveOccurred())

	file, err := ioutil.ReadFile(filepath.Join(dir, "package.json"))
	Expect(err).ToNot(HaveOccurred())
	obj := make(map[string]interface{})
	Expect(json.Unmarshal(file, &obj)).To(Succeed())
	engines, ok := obj["engines"].(map[string]interface{})
	Expect(ok).To(BeTrue())
	engines["node"] = nodejsVersion
	file, err = json.Marshal(obj)
	Expect(err).ToNot(HaveOccurred())
	Expect(ioutil.WriteFile(filepath.Join(dir, "package.json"), file, 0644)).To(Succeed())

	return dir
}

func AddDotProfileScriptToApp(dir string) {
	profilePath := filepath.Join(dir, ".profile")
	Expect(ioutil.WriteFile(profilePath, []byte(`#!/usr/bin/env bash
echo PROFILE_SCRIPT_IS_PRESENT_AND_RAN
`), 0755)).To(Succeed())
}
