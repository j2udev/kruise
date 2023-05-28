package cmd

import (
	"github.com/j2udev/boa"
	"github.com/j2udev/kruise/internal/kruise"
	"github.com/spf13/cobra"
)

func NewDeployCmd() *cobra.Command {
	return boa.NewCmd("deploy").
		WithValidOptions(deployOptions()...).
		WithValidProfiles(deployProfiles()...).
		WithOptionsTemplate().
		WithMinValidArgs(1).
		WithAliases([]string{"dep"}).
		WithShortDescription("Deploy the specified options to your Kubernetes cluster").
		WithRunFunc(deploy).
		WithBoolPFlag("dry-run", "d", false, "output the command being performed under the hood").
		WithBoolPFlag("concurrent", "c", false, "deploy the arguments concurrently (deploys in order based on the 'priority' of each deployment passed)").
		WithBoolPFlag("init", "i", false, "add Helm repositories and create Kubernetes secrets for the specified options").
		Build()
}

func deploy(cmd *cobra.Command, args []string) {
	kruise.Deploy(cmd.Flags(), args)
}

func deployOptions() []boa.Option {
	var opts []boa.Option
	for _, d := range kruise.GetDeployments() {
		args := []string{d.Name}
		args = append(args, d.Aliases...)
		opts = append(opts, boa.Option{Args: args, Desc: d.Description.Deploy})
	}
	return opts
}

func deployProfiles() []boa.Profile {
	var profs []boa.Profile
	for _, p := range kruise.GetDeployProfiles() {
		args := []string{p.Name}
		args = append(args, p.Aliases...)
		profs = append(profs, boa.Profile{Args: args, Opts: p.Items, Desc: p.Description.Deploy})
	}
	return profs
}
