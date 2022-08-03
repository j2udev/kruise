package kruise

import "github.com/spf13/pflag"

func Delete(fs *pflag.FlagSet, args []string) {
	d := GetValidDeployments(args)
	Uninstall(fs, d...)
}

// GetDeleteOptions aggregates delete options from all deployers
func GetDeleteOptions() Options {
	deps := Kfg.Manifest.Deploy.Deployments
	var opts Options
	for k, v := range deps {
		args := []string{k}
		opts = append(opts, NewOption(append(args, v.Aliases...), v.Description.Delete))
	}
	return opts
}

// // GetValidDeleteArgs aggregates valid delete arguments from all deployers
// func GetValidDeleteArgs() []string {
// 	args := GetDeleteOptions().GetValidArgs()
// 	return args
// }

// // GetValidDeployments gets all valid deployments given passed arguments
// func GetValidDeployments(args []string) Installers {
// 	var installers Installers
// 	deps := Kfg.Manifest.Deploy.Deployments
// 	for k, v := range deps {
// 		if contains(args, k) || containsAny(args, v.Aliases...) {
// 			charts := NewHelmDeployment(v.Helm).GetHelmCharts()
// 			manifests := NewKubectlDeployment(v.Kubectl).GetKubectlManifests()
// 			for _, c := range charts {
// 				installers = append(installers, c)
// 			}
// 			for _, m := range manifests {
// 				installers = append(installers, m)
// 			}
// 		}
// 	}
// 	return installers
// }
