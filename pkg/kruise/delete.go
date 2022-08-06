package kruise

import "github.com/spf13/pflag"

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
