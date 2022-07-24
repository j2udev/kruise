package kruise

import (
	"github.com/spf13/cobra"
)

// GetDeleteOptions aggregates delete options from all deployers
func GetDeleteOptions() []Option {
	opts := GetHelmDeleteOptions()
	return opts
}

// GetHelmDeleteOptions gets delete options from the Helm deployer
func GetHelmDeleteOptions() []Option {
	var opts []Option
	deps := NewHelmDeployments(Kfg.Manifest.Delete.Helm)
	for _, dep := range deps {
		opts = append(opts, NewOption(dep.Option))
	}
	return opts
}

// GetValidDeleteArgs aggregates valid delete arguments from all deployers
func GetValidDeleteArgs() []string {
	args := GetValidArgs(GetDeleteOptions())
	return args
}

// Delete is a cobra Run function
func Delete(cmd *cobra.Command, args []string) {
	Uninstall(cmd.Flags(), GetValidDeployments(args)...)
}
