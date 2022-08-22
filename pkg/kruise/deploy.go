package kruise

import (
	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/spf13/pflag"
)

type (
	// Deployment is used to capture the map key as the name field from the
	// latest.Deployment object
	Deployment struct {
		name string
		latest.Deployment
	}
	// Deployments is a slice of Deployment objects
	Deployments []Deployment
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

// GetDeployOptions aggregates deploy options from all deployers
func GetDeployOptions() Options {
	var opts Options
	for k, v := range Kfg.Manifest.Deploy.Deployments {
		args := []string{k}
		opts = append(opts, newOption(append(args, v.Aliases...), v.Description.Deploy))
	}
	return opts
}

// GetDeployProfiles gets deploy profiles from Kruise config
func GetDeployProfiles() Profiles {
	var profs Profiles
	for k, v := range Kfg.Manifest.Deploy.Profiles {
		profs = append(profs, newProfile(k, v))
	}
	return profs
}

// Delete determines valid deployments from args and passes the cobra Cmd
// FlagSet to the Uninstall function
func Delete(fs *pflag.FlagSet, args []string) {
	d := getValidInstallers(args)
	Uninstall(fs, d...)
}

// GetDeleteOptions aggregates delete options from all deployers
func GetDeleteOptions() Options {
	var opts Options
	for k, v := range Kfg.Manifest.Deploy.Deployments {
		args := []string{k}
		opts = append(opts, newOption(append(args, v.Aliases...), v.Description.Delete))
	}
	return opts
}

// GetDeleteProfiles gets deploy profiles from Kruise config
func GetDeleteProfiles() Profiles {
	var profs Profiles
	for k, v := range Kfg.Manifest.Deploy.Profiles {
		p := newProfile(k, v)
		p.Desc = p.Description.Delete
		profs = append(profs, p)
	}
	return profs
}

// GetValidDeployArgs aggregates valid deploy arguments from all deployers
func GetValidDeployArgs() []string {
	args := GetDeployOptions().getValidArgs()
	args = append(args, GetDeployProfiles().getValidArgs()...)
	return args
}

// newDeployment is a helper function for creating a Deployment object with a name
//
// The name is derived from a map entry in a config file and isn't on the
// original latest.Deployment object
func newDeployment(name string, dep latest.Deployment) Deployment {
	return Deployment{name, dep}
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

// getValidInitInstallers gets all valid initial installers given passed arguments
func getValidInitInstallers(args []string) Installers {
	deps := getValidDeployments(args)
	set := make(map[Installer]bool)
	for _, d := range deps {
		secrets := newKubectlDeployment(d.Kubectl).getKubectlSecrets()
		repositories := newHelmDeployment(d.Helm).getHelmRepositories()
		for _, s := range secrets {
			set[s] = true
		}
		for _, r := range repositories {
			set[r] = true
		}
	}
	var installers Installers
	for k := range set {
		installers = append(installers, k)
	}
	return installers
}

// getValidDeployments gets all valid deployments given passed arguments
func getValidDeployments(args []string) map[string]Deployment {
	set := make(map[string]Deployment)
	for _, arg := range args {
		if prof, ok := argIsProfile(arg); ok {
			d := getDeploymentsFromProfile(prof)
			for _, dep := range d {
				set[dep.name] = dep
			}
		} else if dep, ok := argIsDeployment(arg); ok {
			set[dep.name] = dep
		}
	}
	return set
}

// getDeploymentsFromProfile is a helper function for getting all Deployments
// described by a Profile
func getDeploymentsFromProfile(profile Profile) Deployments {
	d := Kfg.Manifest.Deploy.Deployments
	var deps Deployments
	for _, p := range profile.Items {
		deps = append(deps, newDeployment(p, d[p]))
	}
	return deps
}

// argIsDeployment is used to determine if the passed argument is a valid deployment
func argIsDeployment(arg string) (Deployment, bool) {
	d := Kfg.Manifest.Deploy.Deployments
	if dep, ok := d[arg]; ok {
		return newDeployment(arg, dep), true
	}
	for k, v := range d {
		if containsAny(v.Aliases, arg) {
			return newDeployment(k, v), true
		}
	}
	return Deployment{}, false
}

// argIsProfile is used to determine if the passed argument is a valid profile
func argIsProfile(arg string) (Profile, bool) {
	p := Kfg.Manifest.Deploy.Profiles
	if prof, ok := p[arg]; ok {
		return newProfile(arg, prof), true
	}
	for k, v := range p {
		if containsAny(v.Aliases, arg) {
			return newProfile(k, v), true
		}
	}
	return Profile{}, false
}
