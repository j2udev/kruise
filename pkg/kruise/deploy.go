package kruise

import (
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/thoas/go-funk"
)

func GetDeployOptions() []Option {
	opts := GetHelmDeployOptions()
	return opts
}

func GetHelmDeployOptions() []Option {
	var opts []Option
	deps := NewHelmDeployments(Kfg.Manifest.Deploy.Helm)
	for _, dep := range deps {
		opts = append(opts, NewOption(dep.Option))
	}
	return opts
}

func GetValidDeployArgs() []string {
	args := GetValidArgs(GetDeployOptions())
	return args
}

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
		CheckErr(HelmRepoUpdate(shallowDryRun))
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
