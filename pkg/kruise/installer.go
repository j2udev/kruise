package kruise

import (
	"sync"

	"github.com/spf13/pflag"
)

type (
	IInstaller interface {
		Init(dryRun bool) error
		Install(dryRun bool) error
		Uninstall(dryRun bool) error
	}
)

func Install(f *pflag.FlagSet, installer ...IInstaller) {
	shallowDryRun, err := f.GetBool("shallow-dry-run")
	Fatal(err)
	parallel, err := f.GetBool("parallel")
	Fatal(err)
	init, err := f.GetBool("init")
	Fatal(err)
	if init {
		for _, i := range installer {
			if err := i.Init(shallowDryRun); err != nil {
				Fatal(err)
			}
		}
		Fatal(HelmRepoUpdate(shallowDryRun))
	}
	if parallel {
		Logger.Trace("Running in parallel")
		wg := sync.WaitGroup{}
		for _, i := range installer {
			wg.Add(1)
			go installp(i, shallowDryRun, &wg)
		}
		wg.Wait()
		Logger.Trace("Finished running in parallel")
	} else {
		Logger.Trace("Running sequentially")
		for _, i := range installer {
			install(i, shallowDryRun)
		}
		Logger.Trace("Finished running sequentially")
	}
}

func Uninstall(f *pflag.FlagSet, installer ...IInstaller) {
	shallowDryRun, err := f.GetBool("shallow-dry-run")
	Fatal(err)
	parallel, err := f.GetBool("parallel")
	Fatal(err)
	if parallel {
		wg := sync.WaitGroup{}
		for _, i := range installer {
			wg.Add(1)
			go uninstallp(i, shallowDryRun, &wg)
		}
		wg.Wait()
	} else {
		for _, i := range installer {
			uninstall(i, shallowDryRun)
		}
	}
}

func install(i IInstaller, s bool) {
	if err := i.Install(s); err != nil {
		Fatal(err)
	}
}

func installp(i IInstaller, s bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := i.Install(s); err != nil {
		Fatal(err)
	}
}

func uninstall(i IInstaller, s bool) {
	if err := i.Uninstall(s); err != nil {
		Fatal(err)
	}
}

func uninstallp(i IInstaller, s bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := i.Uninstall(s); err != nil {
		Fatal(err)
	}
}
