package helm

import (
	"log"
	"os/exec"

	"github.com/j2udevelopment/kruise/pkg/config"
	"github.com/j2udevelopment/kruise/pkg/utils"
)

// ConstructHelmCommand is used to construct a default Helm install command
// from configuration
func ConstructHelmCommand(helmCmd *config.HelmCommand) {
	if helmCmd.ReleaseName == "" {
		log.Fatal("You must specify a Helm release name")
	}
	if helmCmd.ChartPath == "" {
		log.Fatal("You must specify a Helm chart")
	}
	if helmCmd.Namespace == "" {
		helmCmd.Namespace = "default"
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
// default Helm command
func ConstructAndCheck(helmCmd *config.HelmCommand) {
	CheckHelm()
	ConstructHelmCommand(helmCmd)
}

// DefineInstallArgs applies additional arguments to a default Helm install
// command
func DefineInstallArgs(helmCmd *config.HelmCommand) {
	if len(helmCmd.Args) == 0 {
		helmCmd.Args = []string{
			"helm",
			"upgrade", "-i",
			helmCmd.ReleaseName,
			helmCmd.ChartPath,
			"--namespace", helmCmd.Namespace,
			"--create-namespace",
		}
	}
	if helmCmd.Version != "" {
		helmCmd.Args = append(helmCmd.Args, "--version", helmCmd.Version)
	}
	for _, val := range helmCmd.Values {
		helmCmd.Args = append(helmCmd.Args, "-f", val)
	}
	helmCmd.Args = append(helmCmd.Args, helmCmd.ExtraArgs...)
}

// DefineUninstallArgs applies additional arguments to a default Helm uninstall
// command
func DefineUninstallArgs(helmCmd *config.HelmCommand) {
	if len(helmCmd.Args) == 0 {
		helmCmd.Args = []string{
			"helm",
			"uninstall",
			helmCmd.ReleaseName,
			"--namespace",
			helmCmd.Namespace,
		}
	}
	helmCmd.Args = append(helmCmd.Args, helmCmd.ExtraArgs...)
}

// Install is used to install Helm charts in an abstract way
func Install(shallowDryRun bool, helmCmd *config.HelmCommand) error {
	ConstructAndCheck(helmCmd)
	DefineInstallArgs(helmCmd)
	return utils.ExecuteCommand(shallowDryRun, helmCmd.Args[0], helmCmd.Args[1:]...)
}

// Uninstall is used to uninstall Helm charts in an abstract way
func Uninstall(shallowDryRun bool, helmCmd *config.HelmCommand) error {
	ConstructAndCheck(helmCmd)
	DefineUninstallArgs(helmCmd)
	return utils.ExecuteCommand(shallowDryRun, helmCmd.Args[0], helmCmd.Args[1:]...)
}
