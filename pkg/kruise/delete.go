package kruise

import (
	"sync"

	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/thoas/go-funk"
)

func GetDeleteOptions() []Option {
	opts := GetHelmDeleteOptions()
	return opts
}

func GetHelmDeleteOptions() []Option {
	deps := Kfg.GetDeleteConfig()
	return funk.Reduce(deps.Helm, func(acc []Option, h latest.HelmDeployment) []Option {
		return append(acc, Option{h.Option})
	}, []Option{}).([]Option)
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
