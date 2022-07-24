package kruise

import (
	"sync"

	"github.com/spf13/pflag"
	"github.com/thoas/go-funk"
)

type (
	IInstaller interface {
		Install(dryRun bool) error
		Uninstall(dryRun bool) error
	}
)

func Install(f *pflag.FlagSet, i ...IInstaller) {
	shallowDryRun, err := f.GetBool("shallow-dry-run")
	Fatal(err)
	parallel, err := f.GetBool("parallel")
	Fatal(err)
	init, err := f.GetBool("init")
	Fatal(err)
	if init {
		funk.ForEach(i, func(i IInstaller) {
			if err := i.(HelmDeployment).Init(shallowDryRun); err != nil {
				Fatal(err)
			}
		})
		Fatal(HelmRepoUpdate(shallowDryRun))
	}
	if parallel {
		wg := sync.WaitGroup{}
		funk.ForEach(i, func(i IInstaller) {
			wg.Add(1)
			go func(h IInstaller) {
				defer wg.Done()
				if err := h.Install(shallowDryRun); err != nil {
					Fatal(err)
				}
			}(i)
		})
		wg.Wait()
	} else {
		funk.ForEach(i, func(i IInstaller) {
			if err := i.Install(shallowDryRun); err != nil {
				Fatal(err)
			}
		})
	}
}

func Uninstall(f *pflag.FlagSet, i ...IInstaller) {
	shallowDryRun, err := f.GetBool("shallow-dry-run")
	Fatal(err)
	parallel, err := f.GetBool("parallel")
	Fatal(err)
	if parallel {
		wg := sync.WaitGroup{}
		funk.ForEach(i, func(i IInstaller) {
			wg.Add(1)
			go func(h IInstaller) {
				defer wg.Done()
				if err := h.Uninstall(shallowDryRun); err != nil {
					Fatal(err)
				}
			}(i)
		})
		wg.Wait()
	} else {
		funk.ForEach(i, func(i IInstaller) {
			if err := i.Uninstall(shallowDryRun); err != nil {
				Fatal(err)
			}
		})
	}
}
