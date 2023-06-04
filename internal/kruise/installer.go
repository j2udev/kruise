package kruise

import (
	"sort"
	"sync"

	"github.com/spf13/pflag"
)

type (
	// Installer represents an interface for generic objects that can be
	// installed and uninstalled
	Installer interface {
		Install(fs *pflag.FlagSet)
		Uninstall(fs *pflag.FlagSet)
		GetPriority() int
	}
	// Installers represents a slice of Installer objects
	Installers []Installer
	// KruiseInstaller represents a union type of Kruise Installer implementations
	KruiseInstaller interface {
		HelmRepository | HelmChart | KubectlGenericSecret | KubectlDockerRegistrySecret | KubectlManifest
	}
)

// Init invokes the Install function for all Installers that should only be
// installed during initialization (i.e. HelmRepositories and KubectlSecrets)
func Init(fs *pflag.FlagSet, installers ...Installer) {
	hasHelmDeployment := false
	for _, i := range installers {
		switch d := i.(type) {
		case KubectlDockerRegistrySecret, KubectlGenericSecret:
			i.Install(fs)
		case HelmRepository:
			hasHelmDeployment = true
			i.Install(fs)
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

// installs is used to invoke the install functions of the given Installers
func installs(fs *pflag.FlagSet, installers ...Installer) {
	for _, i := range installers {
		install(i, fs)
	}
}

// installc is used to concurrently invoke the install functions of the given
// Installers
//
// Waitgroups are contructed based on Installer Priority where each batch of
// prioritized deployments are executed concurrently
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
	Logger.Debug("Finished running concurrently")
}

// installp is used to invoke the install function of a given Installer and
// call wg.Done()
func installp(i Installer, fs *pflag.FlagSet, wg *sync.WaitGroup) {
	defer wg.Done()
	install(i, fs)
}

// install is used to invoke the install function of a given Installer
func install(i Installer, fs *pflag.FlagSet) {
	i.Install(fs)
}

// uninstalls is used to invoke the uninstall functions of the given Installers
func uninstalls(fs *pflag.FlagSet, installers ...Installer) {
	for _, i := range installers {
		i.Uninstall(fs)
	}
}

// uninstallc is used to concurrently invoke the uninstall functions of the
// given Installers
//
// Waitgroups are contructed based on Installer Priority where each batch of
// prioritized deployments are executed concurrently
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
	Logger.Debug("Finished running concurrently")
}

// uninstallp is used to invoke the uninstall function of a given Installer and
// call wg.Done()
func uninstallp(i Installer, fs *pflag.FlagSet, wg *sync.WaitGroup) {
	defer wg.Done()
	uninstall(i, fs)
}

// uninstall is used to invoke the uninstall function of a given Installer
func uninstall(i Installer, fs *pflag.FlagSet) {
	i.Uninstall(fs)
}

// priorityMap is used to construct a map of prioritized Installers
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

func toInstallers[T KruiseInstaller](t []T) Installers {
	var installers Installers
	for _, r := range t {
		installers = append(installers, Installer(r))
	}
	return installers
}
