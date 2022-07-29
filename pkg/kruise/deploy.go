package kruise

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

// GetDeployOptions aggregates deploy options from all deployers
func GetDeployOptions() []Option {
	deps := Kfg.Manifest.Deploy.Deployments
	var opts Options
	for k, v := range deps {
		args := []string{k}
		opts = append(opts, NewOption(append(args, v.Aliases...), v.Description.Deploy))
	}
	return opts
}

func GetHelmDeployments() HelmDeployments {
	deps := Kfg.Manifest.Deploy.Deployments
	var h HelmDeployments
	for _, v := range deps {
		h = append(h, HelmDeployment(v.Helm))
	}
	return h
}

// GetValidDeployArgs aggregates valid deploy arguments from all deployers
func GetValidDeployArgs() []string {
	args := GetValidArgs(GetDeployOptions())
	return args
}

// GetValidDeployments gets all valid deployments given passed arguments
func GetValidDeployments(args []string) []IInstaller {
	var validDeployments []IInstaller
	deps := NewHelmDeployments(Kfg.Manifest.Deploy.Helm)
	for _, dep := range deps {
		for _, arg := range args {
			if funk.Contains(strings.Split(dep.Arguments, ", "), arg) {
				validDeployments = append(validDeployments, dep)
			}
		}
	}
	return validDeployments
}

// Deploy is a cobra Run function
func Deploy(cmd *cobra.Command, args []string) {
	Install(cmd.Flags(), GetValidDeployments(args)...)
}

// // GetDeployOptions aggregates deploy options from all deployers
// func GetDeployOptions() []Option {
// 	opts := GetHelmDeployOptions()
// 	return opts
// }

// // GetHelmDeployOptions gets deploy options from the Helm deployer
// func GetHelmDeployOptions() []Option {
// 	var opts []Option
// 	deps := NewHelmDeployments(Kfg.Manifest.Deploy.Helm)
// 	for _, dep := range deps {
// 		opts = append(opts, NewOption(dep.Option))
// 	}
// 	return opts
// }

// // GetValidDeployArgs aggregates valid deploy arguments from all deployers
// func GetValidDeployArgs() []string {
// 	args := GetValidArgs(GetDeployOptions())
// 	return args
// }

// // GetValidDeployments gets all valid deployments given passed arguments
// func GetValidDeployments(args []string) []IInstaller {
// 	var validDeployments []IInstaller
// 	deps := NewHelmDeployments(Kfg.Manifest.Deploy.Helm)
// 	for _, dep := range deps {
// 		for _, arg := range args {
// 			if funk.Contains(strings.Split(dep.Arguments, ", "), arg) {
// 				validDeployments = append(validDeployments, dep)
// 			}
// 		}
// 	}
// 	return validDeployments
// }

// // Deploy is a cobra Run function
// func Deploy(cmd *cobra.Command, args []string) {
// 	Install(cmd.Flags(), GetValidDeployments(args)...)
// }
