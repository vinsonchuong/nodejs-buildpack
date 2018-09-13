package main

import (
	"fmt"
	libbuildpackV3 "github.com/buildpack/libbuildpack"
	"nodejs/v3/detect"
	"os"
)

func main() {
	detectData, err := libbuildpackV3.DefaultDetect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to do default detection: %s", err)
		os.Exit(100)
	}

	err = detect.CreateBuildPlan(&detectData)
	if err != nil {
		detectData.Logger.Debug("failed nodejs detection: %s", err)
		detectData.Fail()
	}

	fmt.Fprint(os.Stdout, detectData.BuildPlan.String())
}
