package cmd

import (
	"github.com/j2udev/boa"
	"github.com/j2udev/kruise/internal/kruise"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	return boa.NewCmd("delete").
		WithValidOptions(deployOptions()...).
		WithValidProfiles(deployProfiles()...).
		WithOptionsTemplate().
		WithMinValidArgs(1).
		WithAliases([]string{"del"}).
		WithShortDescription("Delete the specified options from your Kubernetes cluster").
		WithRunFunc(deploy).
		WithBoolPFlag("dry-run", "d", false, "output the command being performed under the hood").
		WithBoolPFlag("concurrent", "c", false, "delete the arguments concurrently (deletes in order based on the 'priority' of each deployment passed)").
		Build()
}

func delete(cmd *cobra.Command, args []string) {
	kruise.Delete(cmd.Flags(), args)
}

func deleteOptions() []boa.Option {
	var opts []boa.Option
	for _, d := range kruise.GetDeployments() {
		args := []string{d.Name}
		args = append(args, d.Aliases...)
		opts = append(opts, boa.Option{Args: args, Desc: d.Description.Delete})
	}
	return opts
}

func deleteProfiles() []boa.Profile {
	var profs []boa.Profile
	for _, p := range kruise.GetDeployProfiles() {
		args := []string{p.Name}
		args = append(args, p.Aliases...)
		profs = append(profs, boa.Profile{Args: args, Opts: p.Items, Desc: p.Description.Delete})
	}
	return profs
}
