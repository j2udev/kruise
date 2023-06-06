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
	// k8sRef is a helper struct used in deduplicating k8s resources
	k8sRef struct {
		Name      string
		Namespace string
	}
)

// Deploy determines passed deployments from args and passes the cobra Cmd
// FlagSet to the Uninstall function
func Deploy(fs *pflag.FlagSet, args []string) {
	init, err := fs.GetBool("init")
	Fatal(err)
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
	charts := getPassedHelmCharts(args)
	manifests := getPassedKubectlManifests(args)
	repos := getPassedHelmRepos(args)
	secrets := getPassedKubectlSecrets(args)
	installers = append(installers, charts...)
	installers = append(installers, manifests...)
	installers = append(installers, repos...)
	installers = append(installers, secrets...)
	return installers
}

// getPassedKubectlManifests gets all passed KubectlManifests given passed
// arguments
func getPassedKubectlManifests(args []string) Installers {
	deps := getPassedDeployments(args)
	var installers Installers
	manifests := make(map[Installer]bool)
	for _, d := range deps {
		kubectlDeployment := newKubectlDeployment(d.Kubectl)
		man := kubectlDeployment.getKubectlManifests()
		for _, m := range man {
			manifests[m] = true
		}
	}
	for k := range manifests {
		installers = append(installers, k)
	}
	return installers
}

// getPassedKubectlSecrets gets all passed KubectlSecrets given passed arguments
// and consolidates the same secrets that are being created across multiple
// namespaces
func getPassedKubectlSecrets(args []string) Installers {
	deps := getPassedDeployments(args)
	var installers Installers
	gsecrets := make(map[string]KubectlGenericSecret)
	dsecrets := make(map[string]KubectlDockerRegistrySecret)
	for _, d := range deps {
		kubectlDeployment := newKubectlDeployment(d.Kubectl)
		genericSecrets := kubectlDeployment.getKubectlGenericSecrets()
		dockerRegistrySecrets := kubectlDeployment.getKubectlDockerRegistrySecrets()
		for _, s := range genericSecrets {
			hashedSecret := s.hash()
			if val, ok := gsecrets[hashedSecret]; ok {
				val.Namespaces = append(val.Namespaces, s.Namespace)
				gsecrets[hashedSecret] = val
			} else {
				gsecrets[hashedSecret] = s
			}
		}
		for _, s := range dockerRegistrySecrets {
			hashedSecret := s.hash()
			if val, ok := dsecrets[hashedSecret]; ok {
				val.Namespaces = append(val.Namespaces, s.Namespace)
				dsecrets[hashedSecret] = val
			} else {
				dsecrets[hashedSecret] = s
			}
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
func getPassedHelmCharts(args []string) Installers {
	deps := getPassedDeployments(args)
	var installers Installers
	charts := make(map[Installer]bool)
	for _, d := range deps {
		helmDeployment := newHelmDeployment(d.Helm)
		cha := helmDeployment.getHelmCharts()
		for _, c := range cha {
			charts[c] = true
		}
	}
	for k := range charts {
		installers = append(installers, k)
	}
	return installers
}

// getPassedHelmRepos gets all passed HelmRepositories given passed arguments
func getPassedHelmRepos(args []string) Installers {
	deps := getPassedDeployments(args)
	var installers Installers
	repos := make(map[Installer]bool)
	for _, d := range deps {
		helmDeployment := newHelmDeployment(d.Helm)
		repositories := helmDeployment.getHelmRepositories()
		for _, r := range repositories {
			repos[r] = true
		}
	}
	for k := range repos {
		installers = append(installers, k)
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
