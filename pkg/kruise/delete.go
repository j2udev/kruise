package kruise

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/thoas/go-funk"
)

func GetDeleteOptions() []Option {
	opts := GetHelmDeleteOptions()
	return opts
}

func GetHelmDeleteOptions() []Option {
	var opts []Option
	deps := NewHelmDeployments(Kfg.Manifest.Delete.Helm)
	for _, dep := range deps {
		opts = append(opts, NewOption(dep.Option))
	}
	return opts
}

func GetValidDeleteArgs() []string {
	args := GetValidArgs(GetDeleteOptions())
	return args
}

func Uninstall(f *pflag.FlagSet, i ...IInstaller) {
	shallowDryRun, err := f.GetBool("shallow-dry-run")
	CheckErr(err)
	parallel, err := f.GetBool("parallel")
	CheckErr(err)
	if parallel {
		wg := sync.WaitGroup{}
		funk.ForEach(i, func(i IInstaller) {
			wg.Add(1)
			go func(h IInstaller) {
				defer wg.Done()
				if err := h.Uninstall(shallowDryRun); err != nil {
					CheckErr(err)
				}
			}(i)
		})
		wg.Wait()
	} else {
		funk.ForEach(i, func(i IInstaller) {
			if err := i.Uninstall(shallowDryRun); err != nil {
				CheckErr(err)
			}
		})
	}
}

func Delete(cmd *cobra.Command, args []string) {
	Uninstall(cmd.Flags(), GetValidDeployments(args)...)
}
