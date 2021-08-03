package kruise

import (
	"log"
	"os/exec"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// HelmDeployment struct used to unmarshal yaml kruise configuration
type HelmDeployment struct {
	Option      `mapstructure:"option"`
	HelmCommand `mapstructure:"command"`
}

// HelmCommand struct used to unmarshal nested yaml leafctl configuration
type HelmCommand struct {
	ReleaseName        string
	ChartPath          string
	Namespace          string
	Version            string
	Values             []string
	Args               []string
	ExtraInstallArgs   []string
	ExtraUninstallArgs []string
}

// GetHelmDeployments returns a slice of HelmDeployment objects based on
// unmarshalled configuration. Valid keys are `deploy` and `delete`
func GetHelmDeployments(key string) []HelmDeployment {
	var d []HelmDeployment
	cobra.CheckErr(mapstructure.Decode(viper.Get(key+".helm"), &d))
	return d
}

// Install is used to install Helm charts in an abstract way
func (helm *HelmCommand) Install(shallowDryRun bool) error {
	checkHelm()
	constructHelmInstallCommand(helm)
	return ExecuteCommand(shallowDryRun, helm.Args[0], helm.Args[1:]...)
}

// Uninstall is used to install Helm charts in an abstract way
func (helm *HelmCommand) Uninstall(shallowDryRun bool) error {
	checkHelm()
	constructHelmUninstallCommand(helm)
	return ExecuteCommand(shallowDryRun, helm.Args[0], helm.Args[1:]...)
}

// constructHelmInstallCommand is used to construct a default Helm install
// command from configuration
func constructHelmInstallCommand(helm *HelmCommand) {
	if helm.ReleaseName == "" {
		log.Fatal("You must specify a Helm release name")
	}
	if helm.ChartPath == "" {
		log.Fatal("You must specify a Helm chart")
	}
	if helm.Namespace == "" {
		helm.Namespace = "default"
	}
	defineInstallArgs(helm)
}

// constructHelmUninstallCommand is used to construct a default Helm uninstall
// command from configuration
func constructHelmUninstallCommand(helm *HelmCommand) {
	if helm.ReleaseName == "" {
		log.Fatal("You must specify a Helm release name")
	}
	if helm.Namespace == "" {
		helm.Namespace = "default"
	}
	defineUninstallArgs(helm)
}

// checkHelm is used to verify that Helm is installed
func checkHelm() {
	helmCheck := exec.Command("command", "-v", "helm")
	if err := helmCheck.Run(); err != nil {
		log.Fatalf("%s", "Helm does not appear to be installed")
	}
}

// defineInstallArgs applies additional arguments to a default Helm install
// command
func defineInstallArgs(helm *HelmCommand) {
	if len(helm.Args) == 0 {
		helm.Args = []string{
			"helm",
			"upgrade", "-i",
			helm.ReleaseName,
			helm.ChartPath,
			"--namespace", helm.Namespace,
			"--create-namespace",
		}
	}
	if helm.Version != "" {
		helm.Args = append(helm.Args, "--version", helm.Version)
	}
	for _, val := range helm.Values {
		helm.Args = append(helm.Args, "-f", val)
	}
	helm.Args = append(helm.Args, helm.ExtraInstallArgs...)
}

// defineUninstallArgs applies additional arguments to a default Helm uninstall
// command
func defineUninstallArgs(helm *HelmCommand) {
	if len(helm.Args) == 0 {
		helm.Args = []string{
			"helm",
			"uninstall",
			helm.ReleaseName,
			"--namespace",
			helm.Namespace,
		}
	}
	helm.Args = append(helm.Args, helm.ExtraUninstallArgs...)
}
