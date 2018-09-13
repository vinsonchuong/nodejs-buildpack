package detect

import (
	"errors"
	"fmt"
	libbuildpackV3 "github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libbuildpack"
	"os"
	"path/filepath"
)

func CreateBuildPlan(detectData *libbuildpackV3.Detect) (error) {
	packageJSONPath := filepath.Join(detectData.Application.Root, "package.json")
	if exists, err := libbuildpack.FileExists(packageJSONPath); err != nil {
		return fmt.Errorf("error checking filepath %s", packageJSONPath)
	} else if !exists {
		return fmt.Errorf("no package.json found in %s", packageJSONPath)
	}

	pkgJSON, err := loadPackageJSON(packageJSONPath)
	if err != nil {
		return err
	}

	detectData.BuildPlan["node"] = libbuildpackV3.BuildPlanDependency{Version: pkgJSON.Engines.Node}

	return nil
}

type packageJSON struct {
	Engines engines `json:"engines"`
}

type engines struct {
	Node string `json:"node"`
	Yarn string `json:"yarn"`
	NPM  string `json:"npm"`
	Iojs string `json:"iojs"`
}


func loadPackageJSON(path string) (packageJSON, error) {
	var p packageJSON

	err := libbuildpack.NewJSON().Load(path, &p)
	if err != nil && !os.IsNotExist(err) {
		return packageJSON{}, err
	}

	if p.Engines.Node == "" {
		return packageJSON{}, errors.New("node version not specified")
	}

	return p, nil
}