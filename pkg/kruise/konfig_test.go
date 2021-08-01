package kruise

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	defaultKfg := Konfig{
		Metadata: Metadata{
			Path:      "../../test",
			Extension: "yaml",
			Name:      "test-default-kruise-manifest",
		},
	}
	overriddenKfg := Konfig{
		Metadata: Metadata{
			Override: "../../test/test-overridden-kruise-manifest.yaml",
		},
	}
	testDefaultHelmDeployer := []HelmDeployment{
		{
			Option{"jaeger", "test jaeger description"},
			HelmCommand{"jaeger", "test/jaeger", "jaeger", "7.7.7", []string{"values/values.yaml"}, nil, []string{"--wait", "--dry-run"}, nil},
		},
		{
			Option{"kafka", "test kafka description"},
			HelmCommand{"kafka", "test/kafka", "kafka", "8.8.8", []string{"values/values.yaml"}, nil, []string{"--dry-run"}, nil},
		},
		{
			Option{"mongodb, mongo", "test mongodb description"},
			HelmCommand{"mongodb", "test/mongodb", "mongodb", "9.9.9", []string{"values/values.yaml"}, nil, []string{"--reuse-values"}, nil},
		},
	}
	testOverriddenHelmDeployer := []HelmDeployment{
		{
			Option{"jaeger", "test jaeger description"},
			HelmCommand{"jaeger", "test/jaeger", "jaeger", "7.7.7", []string{"values/values.yaml"}, nil, []string{"--wait", "--dry-run"}, nil},
		},
		{
			Option{"kafka", "test kafka description"},
			HelmCommand{"kafka", "test/kafka", "kafka", "8.8.8", []string{"values/values.yaml"}, nil, []string{"--dry-run"}, nil},
		},
	}
	testDefaultHelmDeleter := []HelmDeployment{
		{
			Option{"jaeger", "test jaeger description"},
			HelmCommand{"jaeger", "", "jaeger", "", nil, nil, nil, nil},
		},
		{
			Option{"kafka", "test kafka description"},
			HelmCommand{"kafka", "", "kafka", "", nil, nil, nil, nil},
		},
		{
			Option{"mongodb, mongo", "test mongodb description"},
			HelmCommand{"mongodb", "", "mongodb", "", nil, nil, nil, nil},
		},
	}
	testOverriddenHelmDeleter := []HelmDeployment{
		{
			Option{"jaeger", "test jaeger description"},
			HelmCommand{"jaeger", "", "jaeger", "", nil, nil, nil, nil},
		},
		{
			Option{"kafka", "test kafka description"},
			HelmCommand{"kafka", "", "kafka", "", nil, nil, nil, nil},
		},
	}
	defaultKfg.Initialize()
	for i, dep := range GetHelmDeployments("deploy") {
		assert.Equal(t, testDefaultHelmDeployer[i], dep, "Test configuration did not match unmarshalled/decoded configuration")
	}
	for i, dep := range GetHelmDeployments("delete") {
		assert.Equal(t, testDefaultHelmDeleter[i], dep, "Test configuration did not match unmarshalled/decoded configuration")
	}
	overriddenKfg.Initialize()
	for i, dep := range GetHelmDeployments("deploy") {
		assert.Equal(t, testOverriddenHelmDeployer[i], dep, "Test configuration did not match unmarshalled/decoded configuration")
	}
	for i, dep := range GetHelmDeployments("delete") {
		assert.Equal(t, testOverriddenHelmDeleter[i], dep, "Test configuration did not match unmarshalled/decoded configuration")
	}
}
