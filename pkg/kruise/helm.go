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

// NewHelmRepository is a helper function for dealing with the latest.HelmRepository
// to HelmRepository type definition
func NewHelmRepository(rep latest.HelmRepository) HelmRepository {
	return HelmRepository(rep)
}

// NewHelmRepositories is a helper function for dealing with the latest.HelmRepository
// to HelmRepository type definition
func NewHelmRepositories(reps []latest.HelmRepository) HelmRepositories {
	var r HelmRepositories
	for _, rep := range reps {
		r = append(r, NewHelmRepository(rep))
	}
	return r
}

// NewHelmChart is a helper function for dealing with the latest.HelmChart
// to HelmChart type definition
func NewHelmChart(c latest.HelmChart) HelmChart {
	return HelmChart(c)
}

// NewHelmCharts is a helper function for dealing with the latest.HelmChart
// to HelmChart type definition
func NewHelmCharts(charts []latest.HelmChart) HelmCharts {
	var c HelmCharts
	for _, chart := range charts {
		c = append(c, NewHelmChart(chart))
	}
	return c
}

// GetHelmRepositories is a helper function for grabbing the HelmRepositories
// from a HelmDeployment
func (d HelmDeployment) GetHelmRepositories() HelmRepositories {
	return NewHelmRepositories(d.Repositories)
}

// GetHelmCharts is a helper function for grabbing the HelmCharts
// from a HelmDeployment
func (d HelmDeployment) GetHelmCharts() HelmCharts {
	return NewHelmCharts(d.Charts)
}

func (c HelmChart) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkHelm()
	}
	err = helmExecute(d, c.installArgs(fs))
	if err != nil {
		if strings.Contains(err.Error(), "deprecated") {
			Warn(err)
		} else {
			Error(err)
		}
	}
}

func (c HelmChart) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkHelm()
	}
	helmExecute(d, c.uninstallArgs(fs))
}

func (c HelmChart) GetPriority() int {
	return c.Priority
}

func (r HelmRepository) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkHelm()
	}
	helmExecute(d, r.installArgs(fs))
}

func (r HelmRepository) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkHelm()
	}
	Warn(helmExecute(d, r.uninstallArgs(fs)))
}

func (r HelmRepository) GetPriority() int {
	// For now, HelmRepositories are just installed first
	return 0
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
		u := "***"
		p := []byte("***")
		if !sdr {
			var up, pp string
			up = "Please enter your username for the " + r.Name + " Helm repository"
			pp = "Please enter your password for the " + r.Name + " Helm repository"
			un, pw, err := credentialPrompt(up, pp)
			Fatal(err)
			u = un
			p = []byte(pw)
		}
		args = append(args,
			"--username", u,
			"--password", string(p),
			"--pass-credentials")
	}
	return args
}

func (r HelmRepository) uninstallArgs(fs *pflag.FlagSet) []string {
	Logger.Warn("TODO: helmRepoUninstallArgs")
	return []string{"TODO: HelmRepository.Uninstall()"}
}

func helmRepoUpdate(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	helmExecute(d, []string{"repo", "update"})
}

func helmExecute(dry bool, args []string) error {
	return NewCmd("helm").
		WithArgs(args).
		WithDryRun(dry).
		Build().
		Execute()
}

func checkHelm() {
	err := exec.Command("helm").Run()
	Fatalf(err, "Helm does not appear to be installed")
}
