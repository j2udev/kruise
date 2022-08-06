package kruise

import "github.com/spf13/pflag"

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

// GetValidDeployArgs aggregates valid deploy arguments from all deployers
func GetValidDeployArgs() []string {
	args := GetDeployOptions().getValidArgs()
	return args
}

// getValidDeployments gets all valid deployments given passed arguments
func getValidDeployments(args []string) Installers {
	var installers Installers
	deps := Kfg.Manifest.Deploy.Deployments
	for k, v := range deps {
		if contains(args, k) || containsAny(args, v.Aliases...) {
			charts := NewHelmDeployment(v.Helm).GetHelmCharts()
			manifests := NewKubectlDeployment(v.Kubectl).GetKubectlManifests()
			for _, c := range charts {
				installers = append(installers, c)
			}
			for _, m := range manifests {
				installers = append(installers, m)
			}
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
			secrets := NewKubectlDeployment(v.Kubectl).GetKubectlSecrets()
			repositories := NewHelmDeployment(v.Helm).GetHelmRepositories()
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
