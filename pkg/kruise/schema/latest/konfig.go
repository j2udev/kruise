package latest

import (
	"errors"

	"github.com/j2udevelopment/kruise/pkg/kruise/schema"
	"github.com/thoas/go-funk"
)

type (
	KruiseConfig struct {
		APIVersion string      `mapstructure:"apiVersion"`
		Kind       string      `mapstructure:"kind"`
		Deploy     Deployments `mapstructure:"deploy"`
		Delete     Deployments `mapstructure:"delete"`
	}

	Deployments struct {
		Helm []HelmDeployment `mapstructure:"helm"`
	}

	HelmDeployment struct {
		Option    `mapstructure:"option"`
		HelmChart `mapstructure:"chart"`
		Priority  int `mapstructure:"priority"`
	}

	HelmRepository struct {
		RepoName    string `mapstructure:"name,omitempty"`
		RepoUrl     string `mapstructure:"url,omitempty"`
		PrivateRepo bool   `mapstructure:"private,omitempty"`
	}

	HelmChart struct {
		ChartName     string         `yaml:"chartName"`
		ReleaseName   string         `yaml:"releaseName"`
		Namespace     string         `yaml:"namespace"`
		ChartPath     string         `yaml:"chartPath,omitempty"`
		Repository    HelmRepository `yaml:"repository"`
		Version       string         `yaml:"version"`
		Values        []string       `yaml:"values"`
		SetValues     []string       `yaml:"setValues,omitempty"`
		InstallArgs   []string       `yaml:"installArgs,omitempty"`
		UninstallArgs []string       `yaml:"uninstallArgs,omitempty"`
	}

	Option struct {
		Arguments   string
		Description string
	}
)

func (c KruiseConfig) GetVersion() (string, error) {
	if funk.Contains(schema.Versions, c.APIVersion) {
		return c.APIVersion, nil
	}
	return "", errors.New("invalid APIVersion")
}
