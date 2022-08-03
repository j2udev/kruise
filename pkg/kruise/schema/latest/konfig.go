package latest

import "github.com/j2udevelopment/kruise/pkg/kruise/schema/version"

var Version = "v1alpha2"

func NewKruiseConfig() version.IVersionedConfig {
	return new(KruiseConfig)
}

type (
	KruiseConfig struct {
		APIVersion string       `mapstructure:"apiVersion"`
		Kind       string       `mapstructure:"kind"`
		Deploy     DeployConfig `mapstructure:"deploy"`
	}

	DeployConfig struct {
		Deployments map[string]Deployment `mapstructure:"deployments"`
	}

	Deployment struct {
		Aliases     []string          `mapstructure:"aliases"`
		Description DeploymentDesc    `mapstructure:"description"`
		Helm        HelmDeployment    `mapstructure:"helm"`
		Kubectl     KubectlDeployment `mapstructure:"kubectl"`
	}

	DeploymentDesc struct {
		Deploy string `mapstructure:"deploy"`
		Delete string `mapstructure:"delete"`
	}

	HelmDeployment struct {
		Repositories []HelmRepository `mapstructure:"repositories"`
		Charts       []HelmChart      `mapstructure:"charts"`
	}

	KubectlDeployment struct {
		Secrets   []KubectlSecret   `mapstructure:"secrets"`
		Manifests []KubectlManifest `mapstructure:"manifests"`
	}

	HelmRepository struct {
		Url     string `mapstructure:"url"`
		Name    string `mapstructure:"name"`
		Private bool   `mapstructure:"private"`
	}

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

	KubectlSecret struct {
		Type      string `mapstructure:"type"`
		Name      string `mapstructure:"name"`
		Namespace string `mapstructure:"namespace"`
		Registry  string `mapstructure:"registry,omitempty"`
	}

	KubectlManifest struct {
		Namespace string   `mapstructure:"namespace"`
		Priority  int      `mapstructure:"priority"`
		Paths     []string `mapstructure:"paths"`
	}
)

func (c *KruiseConfig) GetVersion() string {
	return c.APIVersion
}
