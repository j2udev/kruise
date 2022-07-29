package kruise

import (
	"log"
	"os/exec"
	"strings"

	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/spf13/pflag"
)

type (
	HelmDeployment   latest.HelmDeployment
	HelmRepository   latest.HelmRepository
	HelmChart        latest.HelmChart
	HelmDeployments  []HelmDeployment
	HelmRepositories []HelmRepository
	HelmCharts       []HelmChart
	HelmInstallers   interface {
		HelmRepository | HelmChart
	}
)

// NewHelmDeployment is a helper function for dealing with the latest.HelmDeployment
// to HelmDeployment type definition
func NewHelmDeployment(dep latest.HelmDeployment) HelmDeployment {
	return HelmDeployment(dep)
}

// NewHelmDeployments is a helper function for dealing with the latest.HelmDeployment
// to HelmDeployment type definition
func NewHelmDeployments(deps []latest.HelmDeployment) HelmDeployments {
	var d []HelmDeployment
	for _, dep := range deps {
		d = append(d, NewHelmDeployment(dep))
	}
	return d
}

func NewHelmRepository(rep latest.HelmRepository) HelmRepository {
	return HelmRepository(rep)
}

func NewHelmRepositories(reps []latest.HelmRepository) HelmRepositories {
	var r HelmRepositories
	for _, rep := range reps {
		r = append(r, NewHelmRepository(rep))
	}
	return r
}

func NewHelmChart(c latest.HelmChart) HelmChart {
	return HelmChart(c)
}

func NewHelmCharts(charts []latest.HelmChart) HelmCharts {
	var c HelmCharts
	for _, chart := range charts {
		c = append(c, NewHelmChart(chart))
	}
	return c
}

func (d HelmDeployment) GetInstallers() []IInstaller {
	l := len(d.Repositories) + len(d.Charts)
	installers := make([]IInstaller, l)
	for i, r := range d.Repositories {
		installers[i] = IInstaller(NewHelmRepository(r))
	}
	for i, c := range d.Charts {
		installers[i] = IInstaller(NewHelmChart(c))
	}
	return installers
}

func (c HelmChart) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	helmExecute(d, c.installArgs(fs))
}

func (c HelmChart) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	helmExecute(d, c.uninstallArgs(fs))
}

func (r HelmRepository) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	helmExecute(d, r.installArgs(fs))
}

func (r HelmRepository) Uninstall(fs *pflag.FlagSet) {
	Logger.Warn("TODO: HelmRepository.Uninstall()")
}

func (c HelmChart) installArgs(fs *pflag.FlagSet) []string {
	if c.ReleaseName == "" {
		log.Fatal("You must specify a Helm release name")
	}
	if c.ChartPath == "" {
		log.Fatal("You must specify a Helm chart")
	}
	args := []string{
		"upgrade",
		"--install",
		c.ReleaseName,
		c.ChartPath,
		"--namespace",
		c.Namespace,
	}
	if c.Version != "" {
		args = append(args, "--version", c.Version)
	}
	if len(c.Values) > 0 {
		for _, val := range c.Values {
			args = append(args, "-f", val)
		}
	}
	if len(c.SetValues) > 0 {
		for _, val := range c.SetValues {
			args = append(args, "--set", val)
		}
	}
	args = append(args, c.InstallArgs...)
	return args
}

func (c HelmChart) uninstallArgs(fs *pflag.FlagSet) []string {
	if c.ReleaseName == "" {
		log.Fatal("You must specify a Helm release name")
	}
	if c.Namespace == "" {
		c.Namespace = "default"
	}
	args := []string{
		"uninstall",
		c.ReleaseName,
		"--namespace",
		c.Namespace,
	}
	if len(c.UninstallArgs) > 0 {
		args = append(args, c.UninstallArgs...)
	}
	return args
}

func (r HelmRepository) installArgs(fs *pflag.FlagSet) []string {
	sdr, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if r.Name == "" {
		log.Fatal("You must specify a Helm repository name")
	}
	if r.Url == "" {
		log.Fatal("You must specify a Helm repository url")
	}
	args := []string{
		"repo",
		"add",
		r.Name,
		r.Url,
		"--force-update",
	}
	if r.Private {
		usernamePrompt := "Please enter username for " + r.Name + " repository"
		passwordPrompt := "Please enter password for " + r.Name + " repository"
		u, p, err := CredentialPrompt(usernamePrompt, passwordPrompt)
		Fatal(err)
		if sdr {
			u = strings.Repeat("*", len(u))
			p = strings.Repeat("*", len(p))
		}
		args = append(args,
			"--username", u,
			"--password", p,
			"--pass-credentials")
	}
	return args
}

func (r HelmRepository) uninstallArgs(fs *pflag.FlagSet) []string {
	Logger.Warn("TODO: helmRepoUninstallArgs")
	return []string{"TODO: HelmRepository.Uninstall()"}
}

func helmExecute(dry bool, args []string) {
	Fatal(NewCmd("helm").
		WithArgs(args).
		WithDryRun(dry).
		Build().
		Execute())
}

func checkHelm() {
	err := exec.Command("helm").Run()
	Fatalf(err, "Helm does not appear to be installed")
}
