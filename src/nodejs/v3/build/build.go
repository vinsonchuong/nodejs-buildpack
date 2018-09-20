package build

import (
	libbuildpackV3 "github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack"
)

func CreateLaunchMetadata() libbuildpackV3.LaunchMetadata {
	return libbuildpackV3.LaunchMetadata{
		Processes: libbuildpackV3.Processes{
			libbuildpackV3.Process{
				Type:    "web",
				Command: "npm start",
			},
		},
	}
}

func MakeNodeLayers(builder libjavabuildpack.Build) {
}
