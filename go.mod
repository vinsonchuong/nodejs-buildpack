module github.com/cloudfoundry/nodejs-buildpack

require (
	cloud.google.com/go v0.36.0
	github.com/BurntSushi/toml v0.3.1
	github.com/Masterminds/semver v1.4.2
	github.com/blang/semver v3.5.1+incompatible
	github.com/cloudfoundry/cnb-tools v0.0.0-20190219172305-154ea1cc62b2 // indirect
	github.com/cloudfoundry/libbuildpack v0.0.0-20190213200103-30ffb32767ef
	github.com/elazarl/goproxy v0.0.0-20181111060418-2ce16c963a8a
	github.com/golang/mock v1.2.0
	github.com/golang/protobuf v1.2.0
	github.com/google/subcommands v0.0.0-20181012225330-46f0354f6315
	github.com/inconshreveable/go-vhost v0.0.0-20160627193104-06d84117953b
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/rogpeppe/go-charset v0.0.0-20180617210344-2471d30d28b4
	github.com/tidwall/gjson v1.1.3
	github.com/tidwall/match v1.0.1
	go4.org v0.0.0-20190218023631-ce4c26f7be8e
	golang.org/x/build v0.0.0-20190221223049-69dd6b2c22e1
	golang.org/x/crypto v0.0.0-20190219172222-a4c6cb3142f2
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd
	golang.org/x/sys v0.0.0-20190204203706-41f3e6584952
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2
	golang.org/x/tools v0.0.0-20190221204921-83362c3779f5
	gopkg.in/jarcoal/httpmock.v1 v1.0.0-20181117152235-275e9df93516
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/cloudfoundry/libbuildpack => /Users/pivotal/workspace/libbuildpack
