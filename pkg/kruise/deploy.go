package kruise

import (
	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/spf13/pflag"
)

// Deploy determines valid deployments from args and passes the cobra Cmd
// FlagSet to the Uninstall function
func Deploy(fs *pflag.FlagSet, args []string) {
	init, err := fs.GetBool("init")
	Fatal(err)
	d := getValidDeployments(args)
	if init {
		i := getValidInitDeployments(args)
		Init(fs, i...)
	}
	Install(fs, d...)
}

// GetDeployOptions aggregates deploy options from all deployers
func GetDeployOptions() Options {
	deps := Kfg.Manifest.Deploy.Deployments
	var opts Options
	for k, v := range deps {
		args := []string{k}
		opts = append(opts, newOption(append(args, v.Aliases...), v.Description.Deploy))
	}
	return opts
}

// Delete determines valid deployments from args and passes the cobra Cmd
// FlagSet to the Uninstall function
func Delete(fs *pflag.FlagSet, args []string) {
	d := getValidDeployments(args)
	Uninstall(fs, d...)
}

// GetDeleteOptions aggregates delete options from all deployers
func GetDeleteOptions() Options {
	deps := Kfg.Manifest.Deploy.Deployments
	var opts Options
	for k, v := range deps {
		args := []string{k}
		opts = append(opts, newOption(append(args, v.Aliases...), v.Description.Delete))
	}
	return opts
}

// GetValidDeployArgs aggregates valid deploy arguments from all deployers
func GetValidDeployArgs() []string {
	args := GetDeployOptions().getValidArgs()
	return args
}

// getValidDeployments gets all valid deployments given passed arguments
func getValidDeployments(args []string) Installers {
	var installers Installers
	for _, a := range args {
		if d, ok := argIsValid(a); ok {
			charts := newHelmDeployment(d.Helm).getHelmCharts()
			manifests := newKubectlDeployment(d.Kubectl).getKubectlManifests()
			installers = append(installers, toInstallers(charts)...)
			installers = append(installers, toInstallers(manifests)...)
		}
	}
	return installers
}

// getValidInitDeployments gets all valid initial deployments given passed arguments
func getValidInitDeployments(args []string) Installers {
	var x void
	deps := Kfg.Manifest.Deploy.Deployments
	set := make(map[Installer]void)
	for k, v := range deps {
		if contains(args, k) || containsAny(args, v.Aliases...) {
			secrets := newKubectlDeployment(v.Kubectl).getKubectlSecrets()
			repositories := newHelmDeployment(v.Helm).getHelmRepositories()
			for _, s := range secrets {
				set[s] = x
			}
			for _, r := range repositories {
				set[r] = x
			}
		}
	}
	var keys Installers
	for k := range set {
		keys = append(keys, k)
	}
	return keys
}

// argIsValid is used to determine if the passed argument is a valid deployment
// which preserves the order of the passed arguments
func argIsValid(arg string) (latest.Deployment, bool) {
	deps := Kfg.Manifest.Deploy.Deployments
	if _, ok := deps[arg]; ok {
		return deps[arg], true
	}
	for _, v := range deps {
		if containsAny(v.Aliases, arg) {
			return v, true
		}
	}
	return latest.Deployment{}, false
}
