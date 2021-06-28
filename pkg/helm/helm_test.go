package helm

import (
	"strings"
	"testing"

	c "github.com/j2udevelopment/kruise/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestConstructHelmChart(t *testing.T) {
	helmCommands := []string{
		"helm upgrade -i jaeger jaeger/jaeger --namespace jaeger --create-namespace",
		"helm upgrade -i kafka kafka/kafka --namespace kafka --create-namespace",
		"helm upgrade -i mongodb mongodb/mongodb --namespace mongodb --create-namespace",
		"helm upgrade -i mysql mysql/mysql --namespace mysql --create-namespace",
		"helm upgrade -i postgresql postgresql/postgresql --namespace postgresql --create-namespace",
		"helm upgrade -i prom-op prom-op/prom-op --namespace prom-op --create-namespace",
		"helm upgrade -i test!@#$release test!@#$/test!@#$ --namespace test --create-namespace --version 7.7.7 --wait",
	}
	helmConfig := []c.HelmConfig{
		{ReleaseName: "jaeger", ChartPath: "jaeger/jaeger", Namespace: "jaeger"},
		{ReleaseName: "kafka", ChartPath: "kafka/kafka", Namespace: "kafka"},
		{ReleaseName: "mongodb", ChartPath: "mongodb/mongodb", Namespace: "mongodb"},
		{ReleaseName: "mysql", ChartPath: "mysql/mysql", Namespace: "mysql"},
		{ReleaseName: "postgresql", ChartPath: "postgresql/postgresql", Namespace: "postgresql"},
		{ReleaseName: "prom-op", ChartPath: "prom-op/prom-op", Namespace: "prom-op"},
		{ReleaseName: "test!@#$release", ChartPath: "test!@#$/test!@#$", Namespace: "test", Version: "7.7.7", ExtraArgs: []string{"--wait"}},
	}
	assert.Equal(t, len(helmCommands), len(helmConfig), "The length of the helmCommands slice did not match the length of the helmConfig slice")
	for i, config := range helmConfig {
		ConstructHelmCommand(&config)
		DefineInstallArgs(&config)
		// fmt.Println(config)
		assert.Equal(t, helmCommands[i], strings.Join(config.Args, " "), "helmCommand did not match parsed helmConfig")
	}
}
