package dagger

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func BundleBuildpack(sourceDir, destDir string) error {
	if err := copyFile(filepath.Join(sourceDir, "buildpack.toml"), filepath.Join(destDir, "buildpack.toml")); err != nil {
		return err
	}

	if err := os.Mkdir(filepath.Join(destDir, "bin"), os.ModePerm); err != nil {
		return err
	}

	for _, b := range []string{"detect", "build"} {
		cmd := exec.Command(
			"go",
			"build",
			"-o",
			filepath.Join(destDir, "bin", b),
			filepath.Join("nodejs", "v3", b, "cmd"),
		)
		cmd.Env = append(os.Environ(), "GOPATH="+sourceDir, "GOOS=linux")
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
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
