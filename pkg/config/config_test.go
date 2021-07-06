package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalHelmDeployConfig(t *testing.T) {
	cfgFile := File{"../../test", "yaml", "test-kruise-manifest", ""}
	var cfg Manifest
	var helmCfg []HelmDeployment
	testCfgA := HelmDeployment{
		Option{"jaeger", "test jaeger description"},
		HelmCommand{"jaeger", "test/jaeger", "jaeger", "7.7.7", []string{"values/values.yaml"}, nil, []string{"--wait", "--dry-run"}},
	}
	testCfgB := HelmDeployment{
		Option{"kafka", "test kafka description"},
		HelmCommand{"kafka", "test/kafka", "kafka", "8.8.8", []string{"values/values.yaml"}, nil, []string{"--dry-run"}},
	}
	testCfgC := HelmDeployment{
		Option{"mongodb, mongo", "test mongodb description"},
		HelmCommand{"mongodb", "test/mongodb", "mongodb", "9.9.9", []string{"values/values.yaml"}, nil, []string{"--reuse-values"}},
	}
	Initialize(cfgFile, cfg)
	Decode("deploy.helm", &helmCfg)
	assert.Equal(t, testCfgA, helmCfg[0], "Test configuration did not match unmarshalled/decoded configuration")
	assert.Equal(t, testCfgB, helmCfg[1], "Test configuration did not match unmarshalled/decoded configuration")
	assert.Equal(t, testCfgC, helmCfg[2], "Test configuration did not match unmarshalled/decoded configuration")
}

func TestUnmarshalHelmDeleteConfig(t *testing.T) {
	cfgFile := File{"../../test", "yaml", "test-kruise-manifest", ""}
	var cfg Manifest
	var helmCfg []HelmDeployment
	testCfgA := HelmDeployment{
		Option{"jaeger", "test jaeger description"},
		HelmCommand{ReleaseName: "jaeger", Namespace: "jaeger"},
	}
	testCfgB := HelmDeployment{
		Option{"kafka", "test kafka description"},
		HelmCommand{ReleaseName: "kafka", Namespace: "kafka"},
	}
	testCfgC := HelmDeployment{
		Option{"mongodb, mongo", "test mongodb description"},
		HelmCommand{ReleaseName: "mongodb", Namespace: "mongodb"},
	}
	Initialize(cfgFile, cfg)
	Decode("delete.helm", &helmCfg)
	assert.Equal(t, testCfgA, helmCfg[0], "Test configuration did not match unmarshalled/decoded configuration")
	assert.Equal(t, testCfgB, helmCfg[1], "Test configuration did not match unmarshalled/decoded configuration")
	assert.Equal(t, testCfgC, helmCfg[2], "Test configuration did not match unmarshalled/decoded configuration")
}
