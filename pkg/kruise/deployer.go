package kruise

import (
	"strings"

	"github.com/spf13/pflag"
)

// Deployer is used to abstract functionality away from the deploy cobra command
type Deployer struct {
	DeployOptions []Option
	DeleteOptions []Option
	HelmDeployer  []HelmDeployment
	HelmDeleter   []HelmDeployment
}

// NewDeployer is the recommended way to create a Deployer object
func NewDeployer() Deployer {
	return Deployer{
		DeployOptions: GetOptions("deploy"),
		DeleteOptions: GetOptions("delete"),
		HelmDeployer:  GetHelmDeployments("deploy"),
		HelmDeleter:   GetHelmDeployments("delete"),
	}
}

// Deploy checks flags and arguments and deploys to a Kubernetes cluster
// appropriately
func (d Deployer) Deploy(flags *pflag.FlagSet, args []string) {
	shallowDryRun, err := flags.GetBool("shallow-dry-run")
	checkErr(err)
	for _, arg := range args {
		for _, dep := range d.HelmDeployer {
			if contains(strings.Split(dep.Option.Arguments, ", "), arg) {
				err := dep.HelmCommand.Install(shallowDryRun)
				checkErr(err)
			}
		}
	}
}

// Delete checks flags and arguments and deletes from a Kubernetes cluster
// appropriately
func (d Deployer) Delete(flags *pflag.FlagSet, args []string) {
	shallowDryRun, err := flags.GetBool("shallow-dry-run")
	checkErr(err)
	for _, arg := range args {
		for _, dep := range d.HelmDeployer {
			if contains(strings.Split(dep.Option.Arguments, ", "), arg) {
				err := dep.HelmCommand.Uninstall(shallowDryRun)
				checkErr(err)
			}
		}
	}
}
