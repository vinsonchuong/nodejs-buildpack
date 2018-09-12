package v3

import (
	"bytes"
	"fmt"
	"github.com/cloudfoundry/libbuildpack"
	"github.com/cloudfoundry/libbuildpack/cutlass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var _ = Describe("Nodejs V3 buildpack", func() {
	It("should run V3 detection", func() {
		bpDir, err := cutlass.FindRoot()
		Expect(err).ToNot(HaveOccurred())

		workingDir, err := ioutil.TempDir("/tmp", "")
		Expect(err).ToNot(HaveOccurred())
		//defer os.RemoveAll(workingDir)

		fmt.Println(workingDir)

		appDir := filepath.Join(workingDir, "app")
		err = os.Mkdir(appDir, os.ModePerm)
		Expect(err).ToNot(HaveOccurred())

		err = libbuildpack.CopyDirectory(filepath.Join(bpDir, "fixtures", "simple_app"), appDir)
		Expect(err).ToNot(HaveOccurred())

		output := &bytes.Buffer{}
		cmd := exec.Command(
			"docker",
			"run",
			"--rm",
			"-v",
			fmt.Sprintf("%s:/workspace", workingDir),
			"-v",
			fmt.Sprintf("%s:/buildpacks/%s/latest", bpDir, "org.cloudfoundry.buildpacks.nodejs"), // TODO : fix me by getting actual id from buildpack toml
			"bpv3:build",
			"/lifecycle/detector",
			"-order",
			"/buildpacks/org.cloudfoundry.buildpacks.nodejs/latest/fixtures/v3/order.toml",
			"-group",
			"/workspace/group.toml",
			"-plan",
			"/workspace/plan.toml",
		)
		cmd.Stdout = output
		cmd.Stderr = output
		if err = cmd.Run(); err != nil {
			Fail("failed to run V3 detection " + output.String())
		}

		// TODO : check group and plan toml files
	})
})
