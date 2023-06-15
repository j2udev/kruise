package kruise

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/j2udev/kruise/internal/schema/latest"
	"github.com/spf13/pflag"
)

type (
	// HelmDeployment encapsulates Helm objects like HelmRepositories and
	// HelmCharts for a given deployment
	HelmDeployment latest.HelmDeployment
	// HelmRepository represents information about a Helm repository
	HelmRepository latest.HelmRepository
	// HelmChart represents information about a Helm chart
	HelmChart latest.HelmChart
	// HelmDeployments represents a slice of HelmDeployment objects
	HelmDeployments []HelmDeployment
	// HelmRepositories represents a slice of HelmRepository objects
	HelmRepositories []HelmRepository
	// HelmCharts represents a slice of HelmChart objects
	HelmCharts []HelmChart
)

// Install is used to execute a Helm install command
func (c HelmChart) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	if err != nil {
		Logger.Fatal(err)
	}
	if !d {
		checkHelm()
	}
	err = helmExecute(d, c.installArgs(fs))
	if err != nil {
		if strings.Contains(err.Error(), "deprecated") {
			if err != nil {
				Logger.Warn(err)
			}
		} else {
			if err != nil {
				Logger.Error(err)
			}
		}
	}
}

// Install is used to execute a Helm repo add command
func (r HelmRepository) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	if err != nil {
		Logger.Fatal(err)
	}
	if !d {
		checkHelm()
	}
	err = helmExecute(d, r.installArgs(fs))
	if err != nil {
		Logger.Error(err)
	}
}

// Uninstall is used to execute a Helm uninstall command
func (c HelmChart) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	if err != nil {
		Logger.Fatal(err)
	}
	if !d {
		checkHelm()
	}
	err = helmExecute(d, c.uninstallArgs(fs))
	if err != nil {
		Logger.Warn(err)
	}
}

// Uninstall is used to execute a Helm repo remove command
func (r HelmRepository) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	if err != nil {
		Logger.Fatal(err)
	}
	if !d {
		checkHelm()
	}
	err = helmExecute(d, r.uninstallArgs(fs))
	if err != nil {
		Logger.Error(err)
	}
}

// GetPriority is used to get the priority of the installer
func (c HelmChart) GetPriority() int {
	return c.Priority
}

// GetPriority is used to get the priority of the installer
func (r HelmRepository) GetPriority() int {
	// For now, HelmRepositories are just installed first
	return 0
}

// IsInit is used to determine whether the installer should be installed during
// initialization
func (m HelmChart) IsInit() bool {
	return m.Init
}

// IsInit is used to determine whether the installer should be installed during
// initialization
func (m HelmRepository) IsInit() bool {
	return m.Init
}

// newHelmDeployment is a helper function for dealing with the latest.HelmDeployment
// to HelmDeployment type definition
func newHelmDeployment(dep latest.HelmDeployment) HelmDeployment {
	return HelmDeployment(dep)
}

// newHelmRepository is a helper function for dealing with the latest.HelmRepository
// to HelmRepository type definition
func newHelmRepository(rep latest.HelmRepository) HelmRepository {
	return HelmRepository(rep)
}

// newHelmRepositories is a helper function for dealing with the latest.HelmRepository
// to HelmRepository type definition
func newHelmRepositories(reps []latest.HelmRepository) HelmRepositories {
	var r HelmRepositories
	for _, rep := range reps {
		r = append(r, newHelmRepository(rep))
	}
	return r
}

// newHelmChart is a helper function for dealing with the latest.HelmChart
// to HelmChart type definition
func newHelmChart(c latest.HelmChart) HelmChart {
	return HelmChart(c)
}

// newHelmCharts is a helper function for dealing with the latest.HelmChart
// to HelmChart type definition
func newHelmCharts(charts []latest.HelmChart) HelmCharts {
	var c HelmCharts
	for _, chart := range charts {
		c = append(c, newHelmChart(chart))
	}
	return c
}

// getHelmRepositories is a helper function for grabbing the HelmRepositories
// from a HelmDeployment
func (d HelmDeployment) getHelmRepositories() HelmRepositories {
	return newHelmRepositories(d.Repositories)
}

// getHelmCharts is a helper function for grabbing the HelmCharts
// from a HelmDeployment
func (d HelmDeployment) getHelmCharts() HelmCharts {
	return newHelmCharts(d.Charts)
}

// installArgs is used to build Helm install CLI args given a FlagSet
func (c HelmChart) installArgs(fs *pflag.FlagSet) []string {
	if c.ChartName == "" {
		Logger.Fatal("You must specify a Helm chart name")
	}
	if c.RepoName == "" {
		Logger.Fatalf("You must specify a Helm repository for %s", c.ChartName)
	}
	if c.ReleaseName == "" {
		Logger.Fatal("You must specify a Helm release name")
	}
	args := []string{
		"upgrade",
		"--install",
		c.ReleaseName,
		c.RepoName + "/" + c.ChartName,
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

// uninstallArgs is used to build Helm uninstall CLI args given a FlagSet
func (c HelmChart) uninstallArgs(fs *pflag.FlagSet) []string {
	if c.ReleaseName == "" {
		Logger.Fatal("You must specify a Helm release name")
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

// installArgs is used to build Helm repo add CLI args given a FlagSet
func (r HelmRepository) installArgs(fs *pflag.FlagSet) []string {
	d, err := fs.GetBool("dry-run")
	if err != nil {
		Logger.Fatal(err)
	}
	if r.Name == "" {
		Logger.Fatal("You must specify a Helm repository name")
	}
	if r.Url == "" {
		Logger.Fatal("You must specify a Helm repository url")
	}
	args := []string{
		"repo",
		"add",
		r.Name,
		r.Url,
		"--force-update", //TODO: force update as the default behavior is probably overkill; think about adding an override flag or something
	}
	if r.Private {
		u := "***"
		p := "***"
		if !d {
			u = normalInputPrompt(fmt.Sprintf("Please enter your username for the %s Helm repository", r.Name))
			p = sensitiveInputPrompt(fmt.Sprintf("Please enter your password for the %s Helm repository", r.Name))
		}
		args = append(args,
			"--username", u,
			"--password", p,
			"--pass-credentials")
	}
	return args
}

// uninstallArgs is used to build Helm repo remove CLI args given a FlagSet
func (r HelmRepository) uninstallArgs(fs *pflag.FlagSet) []string {
	args := []string{
		"repo",
		"remove",
		r.Name,
	}
	return args
}

// hash is used to facilitate storing HelmCharts in a map
func (c HelmChart) hash() string {
	h := sha1.New()
	h.Write([]byte(c.RepoName))
	h.Write([]byte(c.ChartName))
	h.Write([]byte(c.ReleaseName))
	h.Write([]byte(c.Version))
	for _, v := range c.Values {
		h.Write([]byte(v))
	}
	for _, v := range c.SetValues {
		h.Write([]byte(v))
	}
	for _, v := range c.InstallArgs {
		h.Write([]byte(v))
	}
	for _, v := range c.UninstallArgs {
		h.Write([]byte(v))
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// hash is used to facilitate storing HelmRepositories in a map
func (r HelmRepository) hash() string {
	h := sha1.New()
	h.Write([]byte(r.Name))
	h.Write([]byte(r.Url))
	h.Write([]byte(strconv.FormatBool(r.Private)))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// helmRepoUpdate is used to execute a Helm repo update command
func helmRepoUpdate(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	if err != nil {
		Logger.Fatal(err)
	}
	err = helmExecute(d, []string{"repo", "update"})
	if err != nil {
		Logger.Warn(err)
	}
}

// helmExecute is a helper function for executing a Helm command given a set of
// args; it will print the command instead of executing it if dry is true
func helmExecute(dry bool, args []string) error {
	return NewCmd("helm").
		WithArgs(args).
		WithDryRun(dry).
		Build().
		Execute()
}

// checkHelm is used to determine if the user has the Helm CLI installed
func checkHelm() {
	err := exec.Command("helm").Run()
	if err != nil {
		Logger.Fatal(err)
	}
}
