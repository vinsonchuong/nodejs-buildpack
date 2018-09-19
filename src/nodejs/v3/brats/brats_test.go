package brats_test

import (
	"fmt"
	"github.com/BurntSushi/toml"
	libbuildpackV3 "github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libbuildpack/cutlass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"nodejs/v3/dagger"
	"os"
	"os/exec"
	"path/filepath"
)

var _ = Describe("Nodejs V3 buildpack", func() {
	It("should run V3 detection and build", func() {

		bpSourceDir, err := cutlass.FindRoot()
		Expect(err).ToNot(HaveOccurred())

		tmpBpDir, err := ioutil.TempDir("/tmp", "bp")
		Expect(err).ToNot(HaveOccurred())
		fmt.Println(tmpBpDir)

		err = dagger.BundleBuildpack(bpSourceDir, tmpBpDir)
		Expect(err).ToNot(HaveOccurred())
		//defer os.RemoveAll(tmpBpDir)
		Fail("foo")
		workspaceDir, err := ioutil.TempDir("/tmp", "workspace")
		Expect(err).ToNot(HaveOccurred())
		defer os.RemoveAll(workspaceDir)

		err = os.Chmod(workspaceDir, os.ModePerm)
		Expect(err).ToNot(HaveOccurred())

		appDir := filepath.Join(workspaceDir, "app")
		err = os.Mkdir(appDir, os.ModePerm)
		Expect(err).ToNot(HaveOccurred())

		//err = libbuildpack.CopyDirectory(filepath.Join(bpDir, "fixtures", "simple_app"), appDir)
		//Expect(err).ToNot(HaveOccurred())

		// We must ensure container cannot modify app dir
		//err = os.Chmod(appDir, 0755)
		//Expect(err).ToNot(HaveOccurred())

		// Run detect -----------------------------------------------------------------------------

		cmd := exec.Command(
			"docker",
			"run",
			"--rm",
			"-v",
			fmt.Sprintf("%s:/workspace", workspaceDir),
			"-v",
			fmt.Sprintf("%s:/workspace/app", filepath.Join(bpSourceDir, "fixtures", "simple_app")),
			"-v",
			fmt.Sprintf("%s:/buildpacks/org.cloudfoundry.buildpacks.nodejs/latest", tmpBpDir),
			os.Getenv("CNB_BUILD_IMAGE"),
			"/lifecycle/detector",
			"-buildpacks",
			"/buildpacks",
			"-order",
			"/buildpacks/org.cloudfoundry.buildpacks.nodejs/latest/fixtures/v3/order.toml",
			"-group",
			"/workspace/group.toml",
			"-plan",
			"/workspace/plan.toml",
		)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			Fail("failed to run V3 detection")
		}

		group := struct {
			Buildpacks []struct {
				Id      string `toml:"id"`
				Version string `toml:"version"`
			} `toml:"buildpacks"`
		}{}
		_, err = toml.DecodeFile(filepath.Join(workspaceDir, "group.toml"), &group)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(group.Buildpacks)).To(Equal(1))
		Expect(group.Buildpacks[0].Id).To(Equal("org.cloudfoundry.buildpacks.nodejs"))
		Expect(group.Buildpacks[0].Version).To(Equal("1.6.32"))

		plan := libbuildpackV3.BuildPlan{}
		_, err = toml.DecodeFile(filepath.Join(workspaceDir, "plan.toml"), &plan)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(plan)).To(Equal(1))
		Expect(plan).To(HaveKey("node"))
		Expect(plan["node"].Version).To(Equal("~>10"))

		// Run build -----------------------------------------------------------------------------
		cmd = exec.Command(
			"docker",
			"run",
			"--rm",
			"-v",
			fmt.Sprintf("%s:/workspace", workspaceDir),
			"-v",
			fmt.Sprintf("%s:/workspace/app", filepath.Join(bpSourceDir, "fixtures", "simple_app")),
			"-v",
			fmt.Sprintf("%s:/buildpacks/org.cloudfoundry.buildpacks.nodejs/latest", tmpBpDir),
			"-v",
			fmt.Sprintf("%s:/buildpacks/org.cloudfoundry.buildpacks.nodejs/1.6.32", tmpBpDir),
			os.Getenv("CNB_BUILD_IMAGE"),
			"/lifecycle/builder",
			"-buildpacks",
			"/buildpacks",
			"-group",
			"/workspace/group.toml",
			"-plan",
			"/workspace/plan.toml",
		)

		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			Fail("failed to run V3 build")
		}

		launch := libbuildpackV3.LaunchMetadata{}
		_, err = toml.DecodeFile(filepath.Join(workspaceDir, "org.cloudfoundry.buildpacks.nodejs", "launch.toml"), &launch)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(launch.Processes)).To(Equal(1))
		Expect(launch.Processes[0].Type).To(Equal("web"))
		Expect(launch.Processes[0].Command).To(Equal("npm start"))
	})
})
