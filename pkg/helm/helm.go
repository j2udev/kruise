package helm

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/j2udevelopment/kruise/pkg/config"
)

// ConstructChart function used to initialize Helm chart configuration with default values
func ConstructChart(helmConfig *config.HelmConfig) {
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

// Install function used to install Helm charts in an abstract way
func Install(helmConfig config.HelmConfig) {
	helmCheck := exec.Command("command", "-v", "helm")
	if err := helmCheck.Run(); err != nil {
		log.Fatal("Helm does not appear to be installed")
	}
	ConstructChart(&helmConfig)
	for _, val := range helmConfig.Values {
		helmConfig.Args = append(helmConfig.Args, "-f", val)
	}
	for _, val := range helmConfig.ExtraArgs {
		helmConfig.Args = append(helmConfig.Args, val)
	}
	cmd := exec.Command("helm", helmConfig.Args...)
	fmt.Println(cmd)
	fmt.Println("Attempting to deploy " + helmConfig.ChartPath)
	if err := cmd.Run(); err != nil {
		log.Fatal("Something went wrong, is Kubernetes running?")
	} else {
		fmt.Println("Successfully deployed " + helmConfig.ChartPath)
	}
}

// Uninstall function used to uninstall Helm charts in an abstract way
func Uninstall(helmConfig config.HelmConfig) {
	helmCheck := exec.Command("command", "-v", "helm")
	if err := helmCheck.Run(); err != nil {
		log.Fatal("Helm does not appear to be installed")
	}
	ConstructChart(&helmConfig)
	deleteArgs := []string{
		"uninstall",
		helmConfig.ReleaseName,
		"--namespace",
		helmConfig.Namespace,
	}
	for _, val := range helmConfig.ExtraArgs {
		deleteArgs = append(helmConfig.Args, val)
	}
	cmd := exec.Command("helm", deleteArgs...)
	fmt.Println(cmd)
	fmt.Println("Attempting to delete " + helmConfig.ChartPath)
	if err := cmd.Run(); err != nil {
		log.Fatal("Something went wrong, is Kubernetes running?")
	} else {
		fmt.Println("Successfully deleted " + helmConfig.ChartPath)
	}
}
