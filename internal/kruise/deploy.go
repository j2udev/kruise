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

// getAllPassedInstallers gets all passed deployments given passed arguments
func getAllPassedInstallers(args []string) Installers {
	var installers Installers
	deps := getPassedDeployments(args)
	for _, d := range deps {
		installers = append(installers, getPassedHelmRepos(d)...)
		installers = append(installers, getPassedKubectlSecrets(d)...)
		installers = append(installers, getPassedHelmCharts(d)...)
		installers = append(installers, getPassedKubectlManifests(d)...)
	}
	return installers
}

// getPassedKubectlManifests gets all passed KubectlManifests given passed
// arguments
func getPassedKubectlManifests(deployment Deployment) Installers {
	var installers Installers
	manifestMap := make(map[string]Installer)
	var manifests KubectlManifests
	kubectlDeployment := newKubectlDeployment(deployment.Kubectl)
	man := kubectlDeployment.getKubectlManifests()
	for _, m := range man {
		if _, ok := manifestMap[m.hash()]; !ok {
			manifestMap[m.hash()] = m
			manifests = append(manifests, m)
		}
	}
	for _, v := range manifests {
		installers = append(installers, v)
	}
	return installers
}

// getPassedKubectlSecrets gets all passed KubectlSecrets given passed arguments
// and consolidates the same secrets that are being created across multiple
// namespaces
func getPassedKubectlSecrets(deployment Deployment) Installers {
	var installers Installers
	gsecrets := make(map[string]KubectlGenericSecret)
	dsecrets := make(map[string]KubectlDockerRegistrySecret)
	kubectlDeployment := newKubectlDeployment(deployment.Kubectl)
	genericSecrets := kubectlDeployment.getKubectlGenericSecrets()
	dockerRegistrySecrets := kubectlDeployment.getKubectlDockerRegistrySecrets()
	for _, s := range genericSecrets {
		hashedSecret := s.hash()
		if val, ok := gsecrets[hashedSecret]; ok {
			if !contains[string](val.Namespaces, s.Namespace) {
				val.Namespaces = append(val.Namespaces, s.Namespace)
			}
			gsecrets[hashedSecret] = val
		} else {
			gsecrets[hashedSecret] = s
		}
	}
	for _, s := range dockerRegistrySecrets {
		hashedSecret := s.hash()
		if val, ok := dsecrets[hashedSecret]; ok {
			if !contains[string](val.Namespaces, s.Namespace) {
				val.Namespaces = append(val.Namespaces, s.Namespace)
			}
			dsecrets[hashedSecret] = val
		} else {
			dsecrets[hashedSecret] = s
		}
	}
	for _, v := range gsecrets {
		installers = append(installers, v)
	}
	for _, v := range dsecrets {
		installers = append(installers, v)
	}
	return installers
}

// getPassedHelmCharts gets all passed HelmCharts given passed arguments
func getPassedHelmCharts(deployment Deployment) Installers {
	var installers Installers
	chartMap := make(map[string]Installer)
	var charts HelmCharts
	helmDeployment := newHelmDeployment(deployment.Helm)
	cha := helmDeployment.getHelmCharts()
	for _, c := range cha {
		if _, ok := chartMap[c.hash()]; !ok {
			chartMap[c.hash()] = c
			charts = append(charts, c)
		}
	}
	for _, v := range charts {
		installers = append(installers, v)
	}
	return installers
}

// getPassedHelmRepos gets all passed HelmRepositories given passed arguments
func getPassedHelmRepos(deployment Deployment) Installers {
	var installers Installers
	repoMap := make(map[string]Installer)
	var repos HelmRepositories
	helmDeployment := newHelmDeployment(deployment.Helm)
	repositories := helmDeployment.getHelmRepositories()
	for _, r := range repositories {
		if _, ok := repoMap[r.hash()]; !ok {
			repoMap[r.hash()] = r
			repos = append(repos, r)
		}
	}
	for _, v := range repos {
		installers = append(installers, v)
	}
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
