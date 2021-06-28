package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalDeployConfig(t *testing.T) {
	cfgFile := ConfigFile{"../../test/deploy", "yaml", "deploy-test-config", ""}
	InitCustomConfig(cfgFile)
	cfg := HelmConfig{}
	Decode("deploy.jaeger", &cfg)
	assert.Equal(t, "test-jaeger", cfg.ReleaseName, "The releaseName was not unmarshalled correctly")
	assert.Equal(t, "test/jaeger", cfg.ChartPath, "The chartPath was not unmarshalled correctly")
	assert.Equal(t, "test", cfg.Namespace, "The namespace was not unmarshalled correctly")
	assert.Equal(t, "0.39.5", cfg.Version, "The version was not unmarshalled correctly")
	assert.Equal(t, []string{"values/jaeger-values.yaml"}, cfg.Values, "The values were not unmarshalled correctly")
	assert.Equal(t, []string(nil), cfg.Args, "The args were not unmarshalled correctly")
	assert.Equal(t, []string{"--wait", "--dry-run"}, cfg.ExtraArgs, "The extraArgs were not unmarshalled correctly")
}
