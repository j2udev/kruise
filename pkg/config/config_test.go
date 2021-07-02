package config

import (
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
	assert.Equal(t, []string{"values/values.yaml"}, helmCfg.Values, "The values were not unmarshalled correctly")
	assert.Equal(t, []string(nil), helmCfg.Args, "The args were not unmarshalled correctly")
	assert.Equal(t, []string{"--wait", "--dry-run"}, helmCfg.ExtraArgs, "The extraArgs were not unmarshalled correctly")
}

func TestUnmarshalDynamicDeployConfig(t *testing.T) {
	cfgFile := ConfigFile{"../../test/deploy", "yaml", "dynamic-test-config", ""}
	var cfg DynamicConfig
	var helmCfg []DynamicHelmConfig
	testCfgA := DynamicHelmConfig{
		Option{"jaeger", "test jaeger description"},
		HelmConfig{"jaeger", "test/jaeger", "jaeger", "7.7.7", []string{"values/values.yaml"}, nil, []string{"--wait", "--dry-run"}},
	}
	testCfgB := DynamicHelmConfig{
		Option{"kafka", "test kafka description"},
		HelmConfig{"kafka", "test/kafka", "kafka", "8.8.8", []string{"values/values.yaml"}, nil, []string{"--dry-run"}},
	}
	InitCustomConfig(cfgFile, cfg)
	Decode("deploy.helm", &helmCfg)
	assert.Equal(t, testCfgA, helmCfg[0], "Test configuration did not match unmarshalled/decoded configuration")
	assert.Equal(t, testCfgB, helmCfg[1], "Test configuration did not match unmarshalled/decoded configuration")
}
