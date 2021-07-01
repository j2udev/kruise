package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalDeployConfig(t *testing.T) {
	cfgFile := ConfigFile{"../../test/deploy", "yaml", "deploy-test-config", ""}
	var cfg Config
	var helmCfg HelmConfig
	InitCustomConfig(cfgFile, cfg)
	Decode("deploy.jaeger", &helmCfg)
	assert.Equal(t, "test-jaeger", helmCfg.ReleaseName, "The releaseName was not unmarshalled correctly")
	assert.Equal(t, "test/jaeger", helmCfg.ChartPath, "The chartPath was not unmarshalled correctly")
	assert.Equal(t, "test", helmCfg.Namespace, "The namespace was not unmarshalled correctly")
	assert.Equal(t, "0.39.5", helmCfg.Version, "The version was not unmarshalled correctly")
	assert.Equal(t, []string{"values/jaeger-values.yaml"}, helmCfg.Values, "The values were not unmarshalled correctly")
	assert.Equal(t, []string(nil), helmCfg.Args, "The args were not unmarshalled correctly")
	assert.Equal(t, []string{"--wait", "--dry-run"}, helmCfg.ExtraArgs, "The extraArgs were not unmarshalled correctly")
}

func TestUnmarshalDynamicDeployConfig(t *testing.T) {
	cfgFile := ConfigFile{"../../test/deploy", "yaml", "dynamic-test-config", ""}
	var cfg DynamicConfig
	var helmCfg []DynamicHelmConfig
	InitCustomConfig(cfgFile, cfg)
	Decode("deploy.helm", &helmCfg)
	for _, config := range helmCfg {
		fmt.Println(config.OptionName)
		fmt.Println(config.OptionDescription)
		fmt.Println(config.ReleaseName)
		fmt.Println(config.ChartPath)
		fmt.Println(config.Namespace)
		fmt.Println(config.Version)
	}
}
