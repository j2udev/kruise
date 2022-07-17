package kruise

import (
	"log"
	"os/exec"

	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
)

type (
	HelmDeployment struct {
		latest.HelmDeployment
	}
)

// Init is used to initialize Helm repositories
func (h HelmDeployment) Init(shallowDryRun bool) error {
	if !shallowDryRun {
		checkHelm()
	}
	b := NewCmd("helm").WithArgs(constructHelmInitArgs(h))
	if shallowDryRun {
		b = b.WithDryRun()
	}
	cmd := b.Build()
	return cmd.Execute()
}

// Install is ues to install/upgrade a Helm deployment
func (h HelmDeployment) Install(shallowDryRun bool) error {
	if !shallowDryRun {
		checkHelm()
	}
	b := NewCmd("helm").WithArgs(constructHelmInstallArgs(h))
	if shallowDryRun {
		b = b.WithDryRun()
	}
	cmd := b.Build()
	return cmd.Execute()
}

// Uninstall is ues to uninstall a Helm deployment
func (h HelmDeployment) Uninstall(shallowDryRun bool) error {
	if !shallowDryRun {
		checkHelm()
	}
	b := NewCmd("helm").WithArgs(constructHelmUninstallArgs(h))
	if shallowDryRun {
		b = b.WithDryRun()
	}
	cmd := b.Build()
	return cmd.Execute()
}

func checkHelm() {
	helmCheck := exec.Command("helm")
	if err := helmCheck.Run(); err != nil {
		log.Fatalf("%s", "Helm does not appear to be installed")
	}
}

func constructHelmInitArgs(h HelmDeployment) []string {
	r := h.HelmChart.Repository
	if r.RepoName == "" {
		log.Fatal("You must specify a Helm repository name")
	}
	if r.RepoUrl == "" {
		log.Fatal("You must specify a Helm repository url")
	}
	args := []string{
		"repo",
		"add",
		r.RepoName,
		r.RepoUrl,
	}
	// TODO: deal with private repositories
	return args
}

func constructHelmInstallArgs(h HelmDeployment) []string {
	if h.ReleaseName == "" {
		log.Fatal("You must specify a Helm release name")
	}
	if h.ChartPath == "" {
		log.Fatal("You must specify a Helm chart")
	}
	args := []string{
		"upgrade",
		"--install",
		h.ReleaseName,
		h.ChartPath,
		"--namespace",
		h.Namespace,
		"--create-namespace",
	}
	if h.Version != "" {
		args = append(args, "--version", h.Version)
	}
	if len(h.Values) > 0 {
		for _, val := range h.Values {
			args = append(args, "-f", val)
		}
	}
	if len(h.SetValues) > 0 {
		for _, val := range h.SetValues {
			args = append(args, "--set", val)
		}
	}
	args = append(args, h.InstallArgs...)
	return args
}

func constructHelmUninstallArgs(h HelmDeployment) []string {
	if h.ReleaseName == "" {
		log.Fatal("You must specify a Helm release name")
	}
	if h.Namespace == "" {
		h.Namespace = "default"
	}
	args := []string{
		"uninstall",
		h.ReleaseName,
		"--namespace",
		h.Namespace,
	}
	if len(h.UninstallArgs) > 0 {
		args = append(args, h.UninstallArgs...)
	}
	return args
}
