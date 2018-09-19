package main

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	libbuildpackV3 "github.com/buildpack/libbuildpack"
	"nodejs/v3/build"
	"os"
	"path/filepath"
)

func main() {
	launchDir := os.Args[3]

	builder, err := libbuildpackV3.DefaultBuild()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create default builder: %s", err)
		os.Exit(100)
	}

	launchMetadata := build.CreateLaunchMetadata()

	launchFile, err := os.Create(filepath.Join(launchDir, "launch.toml"))
	if err != nil {
		builder.Logger.Debug("failed to create launch.toml: %s", err)
	}

	encoder := toml.NewEncoder(bufio.NewWriter(launchFile))
	if err := encoder.Encode(launchMetadata); err != nil {
		builder.Logger.Debug("failed to write launch metadata: %s", err)
		builder.Failure(100)
	}
}
