package kruise

import (
	"sort"
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
	concurrent, err := f.GetBool("concurrent")
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
	switch {
	case concurrent:
		Logger.Trace("Running concurrently")
		m := priorityMap(installer...)
		var keys []int
		for k := range m {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			wg := sync.WaitGroup{}
			Logger.Debugf("Priority %d waitgroup starting", k)
			wg.Add(len(m[k]))
			installc(m[k], shallowDryRun, &wg)
			wg.Wait()
			Logger.Debugf("Priority %d waitgroup stopping", k)
		}
		Logger.Trace("Finished running concurrently")
	case parallel:
		Logger.Trace("Running in parallel")
		wg := sync.WaitGroup{}
		for _, i := range installer {
			wg.Add(1)
			go installp(i, shallowDryRun, &wg)
		}
		wg.Wait()
		Logger.Trace("Finished running in parallel")
	default:
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
	concurrent, err := f.GetBool("concurrent")
	Fatal(err)
	parallel, err := f.GetBool("parallel")
	Fatal(err)
	switch {
	case concurrent:
		Logger.Trace("Running concurrently")
		m := priorityMap(installer...)
		var keys []int
		for k := range m {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			wg := sync.WaitGroup{}
			Logger.Debugf("Priority %d waitgroup starting", k)
			wg.Add(len(m[k]))
			uninstallc(m[k], shallowDryRun, &wg)
			wg.Wait()
			Logger.Debugf("Priority %d waitgroup stopping", k)
		}
		Logger.Trace("Finished running concurrently")
	case parallel:
		Logger.Trace("Running in parallel")
		wg := sync.WaitGroup{}
		for _, i := range installer {
			wg.Add(1)
			go uninstallp(i, shallowDryRun, &wg)
		}
		wg.Wait()
		Logger.Trace("Finished running in parallel")
	default:
		Logger.Trace("Running sequentially")
		for _, i := range installer {
			uninstall(i, shallowDryRun)
		}
		Logger.Trace("Finished running sequentially")
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

func installc(installers []IInstaller, s bool, wg *sync.WaitGroup) {
	for _, i := range installers {
		go installp(i, s, wg)
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

func uninstallc(installers []IInstaller, s bool, wg *sync.WaitGroup) {
	for _, i := range installers {
		go uninstallp(i, s, wg)
	}
}

func priorityMap(installers ...IInstaller) map[int][]IInstaller {
	m := make(map[int][]IInstaller)
	for _, i := range installers {
		d := i.(HelmDeployment)
		if val, ok := m[d.Priority]; ok {
			m[d.Priority] = append(val, d)
		} else {
			m[d.Priority] = []IInstaller{d}
		}
	}
	return m
}
