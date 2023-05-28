package latest

import "github.com/j2udev/kruise/internal/kruise/schema/version"

var Version = "v1alpha2"

// NewKruiseConfig represents the schema of the Kruise manifest
func NewKruiseConfig() version.IVersionedConfig {
	return new(KruiseConfig)
}

type (
	// KruiseConfig represents the top level keys of the Kruise manifest file
	KruiseConfig struct {
		APIVersion string       `mapstructure:"apiVersion"`
		Kind       string       `mapstructure:"kind"`
		Deploy     DeployConfig `mapstructure:"deploy"`
	}

	// DeployConfig represents a map of dynamic Deployments
	DeployConfig struct {
		Deployments map[string]Deployment `mapstructure:"deployments"`
		Profiles    map[string]Profile    `mapstructure:"profiles"`
	}

	// Deployment represents a flexible means of mapping multiple Helm and
	// Kubectl installers to a single key
	//
	// Aliases and Description are used to determine how the Deployment appears
	// in the Kruise CLI
	Deployment struct {
		Aliases     []string          `mapstructure:"aliases"`
		Description DeploymentDesc    `mapstructure:"description"`
		Helm        HelmDeployment    `mapstructure:"helm"`
		Kubectl     KubectlDeployment `mapstructure:"kubectl"`
	}

	// Profile represents a flexible means of bundling together other deployments
	Profile struct {
		Aliases     []string       `mapstructure:"aliases"`
		Items       []string       `mapstructure:"items"`
		Description DeploymentDesc `mapstructure:"description"`
	}

	// DeploymentDesc represents the descriptions of the Deployment for the
	// deploy and delete commands
	DeploymentDesc struct {
		Deploy string `mapstructure:"deploy"`
		Delete string `mapstructure:"delete"`
	}

	// HelmDeployment represents multiple Helm repositories and Helm charts
	HelmDeployment struct {
		Repositories []HelmRepository `mapstructure:"repositories"`
		Charts       []HelmChart      `mapstructure:"charts"`
	}

	// KubectlDeployment represents multiple Kubectl secrets and Kubectl manifests
	KubectlDeployment struct {
		Secrets   []KubectlSecret   `mapstructure:"secrets"`
		Manifests []KubectlManifest `mapstructure:"manifests"`
	}

	// HelmRepository represents Helm repository information
	HelmRepository struct {
		Url     string `mapstructure:"url"`
		Name    string `mapstructure:"name"`
		Private bool   `mapstructure:"private"`
	}

	// HelmChart represents Helm chart information
	HelmChart struct {
		ChartName     string   `mapstructure:"chartName"`
		ReleaseName   string   `mapstructure:"releaseName"`
		ChartPath     string   `mapstructure:"chartPath"`
		Namespace     string   `mapstructure:"namespace"`
		Values        []string `mapstructure:"values"`
		SetValues     []string `mapstructure:"setValues"`
		InstallArgs   []string `mapstructure:"installArgs"`
		UninstallArgs []string `mapstructure:"uninstallArgs"`
		Priority      int      `mapstructure:"priority"`
		Version       string   `mapstructure:"version"`
	}

	// KubectlSecret represents Kubectl secret information
	KubectlSecret struct {
		Type      string `mapstructure:"type"`
		Name      string `mapstructure:"name"`
		Namespace string `mapstructure:"namespace"`
		Registry  string `mapstructure:"registry,omitempty"`
	}

	// KubectlManifest represents Kubectl manifest information
	KubectlManifest struct {
		Namespace string   `mapstructure:"namespace"`
		Priority  int      `mapstructure:"priority"`
		Paths     []string `mapstructure:"paths"`
	}
)

// GetVersion is used to get the apiVersion of the Kruise config
func (c *KruiseConfig) GetVersion() string {
	return c.APIVersion
}
