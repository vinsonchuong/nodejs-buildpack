package dagger

// TODO: After finalized, move dagger to its own repo and vendor

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func BundleBuildpack(bpSourceDir, destDir string) error {
	err := copyFile(filepath.Join(bpSourceDir, "buildpack.toml"), filepath.Join(destDir, "buildpack.toml"))
	if err != nil {
		return err
	}
	err = os.Mkdir(filepath.Join(destDir, "bin"), os.ModePerm)
	if err != nil {
		return err
	}

	//export BUILDPACK_DIR=`dirname $(readlink -f ${BASH_SOURCE%/*})`
	//source "$BUILDPACK_DIR/scripts/install_go.sh"

	//GOROOT=$GoInstallDir/go GOPATH=$BUILDPACK_DIR $GoInstallDir/go/bin/go build
	// -o $output_dir/build nodejs/v3/build/cmd

	for _, b := range []string{"detect", "build"} {
		cmd := exec.Command(
			"go",
			"build",
			"-o",
			filepath.Join(destDir, "bin", b),
			filepath.Join("nodejs", "v3", b, "cmd"),
		)
		cmd.Env = append(os.Environ(), "GOPATH="+bpSourceDir)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(from, to string) error {
	source, err := os.Open(from)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(to)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
