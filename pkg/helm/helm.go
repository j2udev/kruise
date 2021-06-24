package helm

import (
	"log"
	"os/exec"

	c "github.com/j2udevelopment/kruise/pkg/config"
	u "github.com/j2udevelopment/kruise/pkg/utils"
)

// ConstructChart function used to initialize Helm chart configuration with default values
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

func CheckHelm() {
	helmCheck := exec.Command("command", "-v", "helm")
	if err := helmCheck.Run(); err != nil {
		log.Fatalf("%s", "Helm does not appear to be installed")
	}
}

func ConstructAndCheck(helmConfig *c.HelmConfig) {
	CheckHelm()
	ConstructChart(helmConfig)
}

// Install function used to install Helm charts in an abstract way
func Install(shallowDryRun bool, helmConfig *c.HelmConfig) {
	ConstructAndCheck(helmConfig)
	for _, val := range helmConfig.Values {
		helmConfig.Args = append(helmConfig.Args, "-f", val)
	}
	for _, val := range helmConfig.ExtraArgs {
		helmConfig.Args = append(helmConfig.Args, val)
	}
	u.ExecuteCommand(shallowDryRun, helmConfig.Args[0], helmConfig.Args[1:]...)
}

// Uninstall function used to uninstall Helm charts in an abstract way
func Uninstall(shallowDryRun bool, helmConfig *c.HelmConfig) {
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
	u.ExecuteCommand(shallowDryRun, deleteArgs[0], deleteArgs[1:]...)
}
