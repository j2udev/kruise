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
	secretRef   struct {
		Name      string
		Namespace string
	}
)

// Deploy determines valid deployments from args and passes the cobra Cmd
// FlagSet to the Uninstall function
func Deploy(fs *pflag.FlagSet, args []string) {
	init, err := fs.GetBool("init")
	Fatal(err)
	d := getValidInstallers(args)
	if init {
		i := getValidInitInstallers(args)
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

// Delete determines valid deployments from args and passes the cobra Cmd
// FlagSet to the Uninstall function
func Delete(fs *pflag.FlagSet, args []string) {
	d := getValidInstallers(args)
	Uninstall(fs, d...)
}

// newDeployment is a helper function for creating a Deployment object from schema
//
// The name is derived from a map entry in a config file and isn't on the
// original latest.Deployment object
func newDeployment(dep latest.Deployment) Deployment {
	return Deployment(dep)
}

// getValidInstallers gets all valid deployments given passed arguments
func getValidInstallers(args []string) Installers {
	deps := getValidDeployments(args)
	var installers Installers
	for _, v := range deps {
		charts := newHelmDeployment(v.Helm).getHelmCharts()
		manifests := newKubectlDeployment(v.Kubectl).getKubectlManifests()
		installers = append(installers, toInstallers(charts)...)
		installers = append(installers, toInstallers(manifests)...)
	}
	return installers
}

// getValidKubectlSecrets gets all valid KubectlSecrets given passed arguments
func getValidKubectlSecrets(args []string) Installers {
	deps := getValidDeployments(args)
	var installers Installers
	secrets := make(map[secretRef]Installer)
	for _, d := range deps {
		kubectlDeployment := newKubectlDeployment(d.Kubectl)
		genericSecrets := kubectlDeployment.getKubectlGenericSecrets()
		dockerRegistrySecrets := kubectlDeployment.getKubectlDockerRegistrySecrets()
		for _, s := range genericSecrets {
			secrets[secretRef{s.Name, s.Namespace}] = s
		}
		for _, s := range dockerRegistrySecrets {
			secrets[secretRef{s.Name, s.Namespace}] = s
		}
	}
	for _, v := range secrets {
		installers = append(installers, v)
	}
	return installers
}

// getValidHelmRepos gets all valid HelmRepositories given passed arguments
func getValidHelmRepos(args []string) Installers {
	deps := getValidDeployments(args)
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

// getValidInitInstallers gets all valid initial installers given passed arguments
func getValidInitInstallers(args []string) Installers {
	var installers Installers
	installers = append(installers, getValidKubectlSecrets(args)...)
	installers = append(installers, getValidHelmRepos(args)...)
	return installers
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

// getValidDeployments gets all valid deployments given passed arguments
// func getValidDeployments(args []string) map[string]Deployment {
func getValidDeployments(args []string) Deployments {
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

// argIsDeployment is used to determine if the passed argument is a valid deployment
func argIsDeployment(arg string) (Deployment, bool) {
	deployments := Kfg.Manifest.Deploy.Deployments
	for _, dep := range deployments {
		if dep.Name == arg || contains(dep.Aliases, arg) {
			return newDeployment(dep), true
		}
	}
	return Deployment{}, false
}

// argIsProfile is used to determine if the passed argument is a valid profile
func argIsProfile(arg string) (Profile, bool) {
	profiles := Kfg.Manifest.Deploy.Profiles
	for _, prof := range profiles {
		if prof.Name == arg || contains(prof.Aliases, arg) {
			return newProfile(prof), true
		}
	}
	return Profile{}, false
}
