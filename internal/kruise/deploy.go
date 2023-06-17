package kruise

import (
	"github.com/j2udev/kruise/internal/schema/latest"
	"github.com/spf13/pflag"
)

type (
	// Deployment is used to capture the map key as the name field from the
	// latest.Deployment object
	Deployment latest.Deployment
	// Deployments is a slice of Deployment objects
	Deployments []Deployment
)

// Deploy determines passed deployments from args and passes the cobra Cmd
// FlagSet to the Uninstall function
func Deploy(fs *pflag.FlagSet, args []string) {
	init, err := fs.GetBool("init")
	if err != nil {
		Logger.Fatal(err)
	}
	d := getPassedInstallers(args)
	if init {
		i := getPassedInitInstallers(args)
		Init(fs, i...)
	}
	Install(fs, d...)
}

// GetDeployments gets deployments from Kruise config
func GetDeployments() Deployments {
	var deps Deployments
	for _, d := range Kfg.Manifest.Deploy.Deployments {
		deps = append(deps, newDeployment(d))
	}
	return deps
}

// GetDeployProfiles gets deploy profiles from Kruise config
func GetDeployProfiles() Profiles {
	var profs Profiles
	for _, p := range Kfg.Manifest.Deploy.Profiles {
		profs = append(profs, newProfile(p))
	}
	return profs
}

// Delete determines passed deployments from args and passes the cobra Cmd
// FlagSet to the Uninstall function
func Delete(fs *pflag.FlagSet, args []string) {
	d := getPassedInstallers(args)
	Uninstall(fs, d...)
}

// newDeployment is a helper function for creating a Deployment object from schema
//
// The name is derived from a map entry in a config file and isn't on the
// original latest.Deployment object
func newDeployment(dep latest.Deployment) Deployment {
	return Deployment(dep)
}

// getPassedInstallers gets all passed deployments given passed arguments
func getPassedInstallers(args []string) Installers {
	allInstallers := getAllPassedInstallers(args)
	var installers Installers
	for _, i := range allInstallers {
		if !i.IsInit() {
			installers = append(installers, i)
		}
	}
	return installers
}

// getPassedInitInstallers gets all passed initial installers given passed arguments
func getPassedInitInstallers(args []string) Installers {
	allInstallers := getAllPassedInstallers(args)
	var installers Installers
	for _, i := range allInstallers {
		if i.IsInit() {
			installers = append(installers, i)
		}
	}
	return installers
}

// getAllPassedInstallers gets all passed installers given passed arguments
func getAllPassedInstallers(args []string) Installers {
	deps := getPassedDeployments(args)
	var installers Installers
	var preInstallers Installers
	var postInstallers Installers
	repoMap := make(map[string]Installer)
	gsecretMap := make(map[string]KubectlGenericSecret)
	dsecretMap := make(map[string]KubectlDockerRegistrySecret)
	chartMap := make(map[string]Installer)
	manifestMap := make(map[string]Installer)
	for _, d := range deps {
		helmDeployment := newHelmDeployment(d.Helm)
		kubectlDeployment := newKubectlDeployment(d.Kubectl)
		repositories := helmDeployment.getHelmRepositories()
		genericSecrets := kubectlDeployment.getKubectlGenericSecrets()
		dockerRegistrySecrets := kubectlDeployment.getKubectlDockerRegistrySecrets()
		cha := helmDeployment.getHelmCharts()
		man := kubectlDeployment.getKubectlManifests()
		for _, r := range repositories {
			if _, ok := repoMap[r.hash()]; !ok {
				repoMap[r.hash()] = r
				preInstallers = append(preInstallers, repoMap[r.hash()])
			}
		}
		for _, s := range genericSecrets {
			hashedSecret := s.hash()
			if val, ok := gsecretMap[hashedSecret]; ok {
				val.Namespaces = append(val.Namespaces, s.Namespace)
				gsecretMap[hashedSecret] = val
			} else {
				gsecretMap[hashedSecret] = s
			}
		}
		for _, s := range dockerRegistrySecrets {
			hashedSecret := s.hash()
			if val, ok := dsecretMap[hashedSecret]; ok {
				val.Namespaces = append(val.Namespaces, s.Namespace)
				dsecretMap[hashedSecret] = val
			} else {
				dsecretMap[hashedSecret] = s
			}
		}
		for _, c := range cha {
			if _, ok := chartMap[c.hash()]; !ok {
				chartMap[c.hash()] = c
				postInstallers = append(postInstallers, c)
			}
		}
		for _, m := range man {
			if _, ok := manifestMap[m.hash()]; !ok {
				manifestMap[m.hash()] = m
				postInstallers = append(postInstallers, m)
			}
		}
	}
	for _, v := range gsecretMap {
		preInstallers = append(preInstallers, v)
	}
	for _, v := range dsecretMap {
		preInstallers = append(preInstallers, v)
	}
	installers = append(installers, preInstallers...)
	installers = append(installers, postInstallers...)
	return installers
}

// getPassedDeployments gets all passed deployments given passed arguments
// func getPassedDeployments(args []string) map[string]Deployment {
func getPassedDeployments(args []string) Deployments {
	deployments := Kfg.Manifest.Deploy.Deployments
	var deps Deployments
	dedup := deduplicateArgs(args)
	for _, arg := range dedup {
		for _, dep := range deployments {
			if dep.Name == arg || contains(dep.Aliases, arg) {
				deps = append(deps, newDeployment(dep))
			}
		}
	}
	return deps
}

// deduplicateArgs is used to deduplicate the given args
//
// It breaks down any profiles into their respective items and does not use a
// set in order to preserve order
func deduplicateArgs(args []string) []string {
	var dedup []string
	for _, a := range args {
		if p, ok := argIsProfile(a); ok {
			if !containsAny(dedup, p.Items...) {
				dedup = append(dedup, p.Items...)
			} else {
				for _, item := range p.Items {
					if !contains(dedup, item) {
						dedup = append(dedup, item)
					}
				}
			}
		} else if dep, ok := argIsDeployment(a); ok {
			if !contains(dedup, dep.Name) {
				dedup = append(dedup, dep.Name)
			}
		}
	}
	return dedup
}

// argIsDeployment is used to determine if the passed argument is a passed deployment
func argIsDeployment(arg string) (Deployment, bool) {
	deployments := Kfg.Manifest.Deploy.Deployments
	for _, dep := range deployments {
		if dep.Name == arg || contains(dep.Aliases, arg) {
			return newDeployment(dep), true
		}
	}
	return Deployment{}, false
}

// argIsProfile is used to determine if the passed argument is a passed profile
func argIsProfile(arg string) (Profile, bool) {
	profiles := Kfg.Manifest.Deploy.Profiles
	for _, prof := range profiles {
		if prof.Name == arg || contains(prof.Aliases, arg) {
			return newProfile(prof), true
		}
	}
	return Profile{}, false
}
