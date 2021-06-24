package helm

import (
	"log"
	"os/exec"

	c "github.com/j2udevelopment/kruise/pkg/config"
	u "github.com/j2udevelopment/kruise/pkg/utils"
)

// ConstructChart is used to construct a default Helm chart from configuration
func ConstructChart(helmConfig *c.HelmConfig) {
	if helmConfig.ReleaseName == "" {
		log.Fatal("You must specify a Helm release name")
	}
	if helmConfig.ChartPath == "" {
		log.Fatal("You must specify a Helm chart")
	}
	if helmConfig.Namespace == "" {
		helmConfig.Namespace = "default"
	}
	if len(helmConfig.Args) == 0 {
		helmConfig.Args = []string{
			"helm",
			"upgrade", "-i",
			helmConfig.ReleaseName,
			helmConfig.ChartPath,
			"--namespace", helmConfig.Namespace,
			"--create-namespace",
		}
	}
	if helmConfig.Version != "" {
		helmConfig.Args = append(helmConfig.Args, "--version", helmConfig.Version)
	}
}

// CheckHelm is used to verify that Helm is installed
func CheckHelm() {
	helmCheck := exec.Command("command", "-v", "helm")
	if err := helmCheck.Run(); err != nil {
		log.Fatalf("%s", "Helm does not appear to be installed")
	}
}

// ConstructAndCheck is used to verify that Helm is installed and consruct a
// default Helm chart from configuration
func ConstructAndCheck(helmConfig *c.HelmConfig) {
	CheckHelm()
	ConstructChart(helmConfig)
}

// Install is used to install Helm charts in an abstract way
func Install(shallowDryRun bool, helmConfig *c.HelmConfig) error {
	ConstructAndCheck(helmConfig)
	for _, val := range helmConfig.Values {
		helmConfig.Args = append(helmConfig.Args, "-f", val)
	}
	for _, val := range helmConfig.ExtraArgs {
		helmConfig.Args = append(helmConfig.Args, val)
	}
	return u.ExecuteCommand(shallowDryRun, helmConfig.Args[0], helmConfig.Args[1:]...)
}

// Uninstall is used to uninstall Helm charts in an abstract way
func Uninstall(shallowDryRun bool, helmConfig *c.HelmConfig) error {
	ConstructAndCheck(helmConfig)
	deleteArgs := []string{
		"helm",
		"uninstall",
		helmConfig.ReleaseName,
		"--namespace",
		helmConfig.Namespace,
	}
	for _, val := range helmConfig.ExtraArgs {
		deleteArgs = append(helmConfig.Args, val)
	}
	return u.ExecuteCommand(shallowDryRun, deleteArgs[0], deleteArgs[1:]...)
}
