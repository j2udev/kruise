package kruise

import "github.com/spf13/pflag"

func Deploy(fs *pflag.FlagSet, args []string) {
	// sdr, err := fs.GetBool("shallow-dry-run")
	// Fatal(err)
	// init, err := fs.GetBool("init")
	// Fatal(err)
	// if !sdr {
	// 	kubeconfig = GetKubeconfig()
	// }
	d := GetValidDeployments(args)
	// if init {
	// 	i := GetValidInitDeployments(args)
	// 	Init(fs, i)
	// }
	Install(fs, d...)
}

// GetDeployOptions aggregates deploy options from all deployers
func GetDeployOptions() Options {
	deps := Kfg.Manifest.Deploy.Deployments
	var opts Options
	for k, v := range deps {
		args := []string{k}
		opts = append(opts, NewOption(append(args, v.Aliases...), v.Description.Deploy))
	}
	return opts
}

// GetValidDeployArgs aggregates valid deploy arguments from all deployers
func GetValidDeployArgs() []string {
	args := GetDeployOptions().GetValidArgs()
	return args
}

// GetValidDeployments gets all valid deployments given passed arguments
func GetValidDeployments(args []string) Installers {
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
