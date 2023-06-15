package latest

import "github.com/j2udev/kruise/internal/schema/version"

var Version = "v1alpha3"

// NewKruiseConfig represents the schema of the Kruise manifest
func NewKruiseConfig() version.IVersionedConfig {
	return new(KruiseConfig)
}

type (
	// KruiseConfig represents the top level keys of the Kruise manifest file
	KruiseConfig struct {
		APIVersion string       `mapstructure:"apiVersion"`
		Kind       string       `mapstructure:"kind"`
		Logger     LoggerConfig `mapstructure:"logger"`
		Deploy     DeployConfig `mapstructure:"deploy"`
	}

	// LoggerConfig is used to define charm log configuration
	LoggerConfig struct {
		Level      string `mapstructure:"level"`
		Caller     bool   `mapstructure:"enableCaller"`
		TimeStamp  bool   `mapstructure:"enableTimestamp"`
		TimeFormat string `mapstructure:"timeFormat"`
	}

	// DeployConfig represents a map of dynamic Deployments
	DeployConfig struct {
		Deployments []Deployment `mapstructure:"deployments"`
		Profiles    []Profile    `mapstructure:"profiles"`
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
		Name        string            `mapstructure:"name"`
	}

	// Profile represents a flexible means of bundling together other deployments
	Profile struct {
		Aliases     []string       `mapstructure:"aliases"`
		Items       []string       `mapstructure:"items"`
		Name        string         `mapstructure:"name"`
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
		Secrets   KubectlSecrets    `mapstructure:"secrets"`
		Manifests []KubectlManifest `mapstructure:"manifests"`
	}

	// HelmRepository represents Helm repository information
	HelmRepository struct {
		Url     string `mapstructure:"url"`
		Name    string `mapstructure:"name"`
		Private bool   `mapstructure:"private"`
		Init    bool   `mapstructure:"init"`
	}

	// HelmChart represents Helm chart information
	HelmChart struct {
		ChartName     string   `mapstructure:"chartName"`
		ReleaseName   string   `mapstructure:"releaseName"`
		RepoName      string   `mapstructure:"repoName"`
		Namespace     string   `mapstructure:"namespace"`
		Values        []string `mapstructure:"values"`
		SetValues     []string `mapstructure:"setValues"`
		InstallArgs   []string `mapstructure:"installArgs"`
		UninstallArgs []string `mapstructure:"uninstallArgs"`
		Priority      int      `mapstructure:"priority"`
		Version       string   `mapstructure:"version"`
		Init          bool     `mapstructure:"init"`
	}

	// KubectlSecrets represents different types of Kubernetes secrets
	KubectlSecrets struct {
		Generic        []KubectlGenericSecret        `mapstructure:"generic"`
		DockerRegistry []KubectlDockerRegistrySecret `mapstructure:"dockerRegistry"`
	}

	// KubectlGenericSecret represents a generic Kubernetes secret
	KubectlGenericSecret struct {
		Name      string   `mapstructure:"name"`
		Namespace string   `mapstructure:"namespace"`
		Literal   []KeyVal `mapstructure:"literal"`
		Init      bool     `mapstructure:"init"`
	}

	// KubectlDockerRegistrySecret represents a docker-registry Kubernetes secret
	KubectlDockerRegistrySecret struct {
		Name      string `mapstructure:"name"`
		Namespace string `mapstructure:"namespace"`
		Registry  string `mapstructure:"registry"`
		Init      bool   `mapstructure:"init"`
	}

	// KeyVal is used to defined key values pairs as separate parameters
	KeyVal struct {
		Key string `mapstructure:"key"`
		Val string `mapstructure:"value"`
	}

	// KubectlManifest represents Kubectl manifest information
	KubectlManifest struct {
		Namespace string   `mapstructure:"namespace"`
		Priority  int      `mapstructure:"priority"`
		Paths     []string `mapstructure:"paths"`
		Init      bool     `mapstructure:"init"`
	}
)

// GetVersion is used to get the apiVersion of the Kruise config
func (c *KruiseConfig) GetVersion() string {
	return c.APIVersion
}
