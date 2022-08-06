package kruise

import (
	"sort"
	"sync"

	"github.com/spf13/pflag"
)

type (
	Installer interface {
		Install(fs *pflag.FlagSet)
		Uninstall(fs *pflag.FlagSet)
		GetPriority() int
	}
	Installers []Installer
)

// Init invokes the Install function for all Installers that should only be
// installed during initialization (i.e. HelmRepositories and KubectlSecrets)
func Init(fs *pflag.FlagSet, installers ...Installer) {
	hasHelmDeployment := false
	u := make(map[string]Installer)
	for _, i := range installers {
		switch d := i.(type) {
		case KubectlSecret:
			// TODO: deduplicate the creation of secrets
			i.Install(fs)
		case HelmRepository:
			hasHelmDeployment = true
			repo := d.Url
			if _, ok := u[repo]; !ok {
				u[repo] = i
				i.Install(fs)
			}
		default:
			Logger.Errorf("Invalid installer for the Init() function: %v", d)
		}
	}
	// if a Helm installer was in the list of the installers to initialize,
	// perform a helm repo update at the end
	if hasHelmDeployment {
		helmRepoUpdate(fs)
	}
}

// Install invokes the Install function for all Installers passed
func Install(fs *pflag.FlagSet, installers ...Installer) {
	concurrent, err := fs.GetBool("concurrent")
	Fatal(err)
	switch {
	case concurrent:
		installc(fs, installers...)
	default:
		installs(fs, installers...)
	}
}

// Uninstall invokes the Uninstall function for all Installers passed
func Uninstall(fs *pflag.FlagSet, installers ...Installer) {
	concurrent, err := fs.GetBool("concurrent")
	Fatal(err)
	switch {
	case concurrent:
		uninstallc(fs, installers...)
	default:
		uninstalls(fs, installers...)
	}
}

func installs(fs *pflag.FlagSet, installers ...Installer) {
	for _, i := range installers {
		install(i, fs)
	}
}

func installc(fs *pflag.FlagSet, installers ...Installer) {
	m := priorityMap(installers...)
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		wg := sync.WaitGroup{}
		Logger.Infof("Priority %d waitgroup starting", k)
		wg.Add(len(m[k]))
		for _, i := range m[k] {
			go installp(i, fs, &wg)
		}
		wg.Wait()
		Logger.Debugf("Priority %d waitgroup stopping", k)
	}
	Logger.Trace("Finished running concurrently")
}

func installp(i Installer, fs *pflag.FlagSet, wg *sync.WaitGroup) {
	defer wg.Done()
	install(i, fs)
}

func install(i Installer, fs *pflag.FlagSet) {
	i.Install(fs)
}

func uninstalls(fs *pflag.FlagSet, installers ...Installer) {
	for _, i := range installers {
		i.Uninstall(fs)
	}
}

func uninstallc(fs *pflag.FlagSet, installers ...Installer) {
	m := priorityMap(installers...)
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		wg := sync.WaitGroup{}
		Logger.Infof("Priority %d waitgroup starting", k)
		wg.Add(len(m[k]))
		for _, i := range m[k] {
			go uninstallp(i, fs, &wg)
		}
		wg.Wait()
		Logger.Debugf("Priority %d waitgroup stopping", k)
	}
	Logger.Trace("Finished running concurrently")
}

func uninstallp(i Installer, fs *pflag.FlagSet, wg *sync.WaitGroup) {
	defer wg.Done()
	uninstall(i, fs)
}

func uninstall(i Installer, fs *pflag.FlagSet) {
	i.Uninstall(fs)
}

func priorityMap(installers ...Installer) map[int]Installers {
	m := make(map[int]Installers)
	for _, installer := range installers {
		p := installer.GetPriority()
		if existingEntry, ok := m[p]; ok {
			m[p] = append(existingEntry, installer)
		} else {
			m[p] = Installers{installer}
		}
	}
	return m
}
