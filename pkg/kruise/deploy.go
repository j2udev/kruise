package kruise

import (
	"strings"
	"sync"

	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/thoas/go-funk"
)

func GetDeployOptions() []Option {
	opts := GetHelmDeployOptions()
	return opts
}

func GetHelmDeployOptions() []Option {
	deps := Kfg.GetDeployConfig()
	return funk.Reduce(deps.Helm, func(acc []Option, h latest.HelmDeployment) []Option {
		return append(acc, Option{h.Option})
	}, []Option{}).([]Option)
}

func GetHelmDeployments() []HelmDeployment {
	deps := Kfg.GetDeployConfig()
	return funk.Map(deps.Helm, func(h latest.HelmDeployment) HelmDeployment {
		return HelmDeployment{h}
	}).([]HelmDeployment)
}

func GetValidDeployments(args []string) []IInstaller {
	h := GetHelmDeployments()
	var validDeployments []IInstaller
	for _, dep := range h {
		for _, arg := range args {
			if funk.Contains(strings.Split(dep.Arguments, ", "), arg) {
				validDeployments = append(validDeployments, dep)
			}
		}
	}
	return validDeployments
}

func Install(f *pflag.FlagSet, i ...IInstaller) {
	shallowDryRun, err := f.GetBool("shallow-dry-run")
	CheckErr(err)
	parallel, err := f.GetBool("parallel")
	CheckErr(err)
	init, err := f.GetBool("init")
	CheckErr(err)
	if init {
		funk.ForEach(i, func(i IInstaller) {
			if err := i.(HelmDeployment).Init(shallowDryRun); err != nil {
				CheckErr(err)
			}
		})
		HelmRepoUpdate(shallowDryRun)
	}
	if parallel {
		wg := sync.WaitGroup{}
		funk.ForEach(i, func(i IInstaller) {
			wg.Add(1)
			go func(h IInstaller) {
				defer wg.Done()
				if err := h.Install(shallowDryRun); err != nil {
					CheckErr(err)
				}
			}(i)
		})
		wg.Wait()
	} else {
		funk.ForEach(i, func(i IInstaller) {
			if err := i.Install(shallowDryRun); err != nil {
				CheckErr(err)
			}
		})
	}
}

func Deploy(cmd *cobra.Command, args []string) {
	Install(cmd.Flags(), GetValidDeployments(args)...)
}
