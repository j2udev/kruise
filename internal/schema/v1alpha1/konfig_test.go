package v1alpha1

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	validConfig = []byte(`
apiVersion: v1alpha1
kind: Config
deploy:
  helm:
    - option:
        arguments: "jaeger"
        description: "Deploys Jaeger to your Kubernetes cluster"
      priority: 1
      chart:
        repository:
          url: "https://jaegertracing.github.io/helm-charts"
          name: jaegertracing
          private: false
        chartName: jaeger
        releaseName: jaeger
        chartPath: jaegertracing/jaeger
        namespace: observability
        version: 0.39.5
        values:
          - /path/to/values/jaeger/jaeger-values.yaml
        setValues:
          - "image=test:0.1.0"
        installArgs:
        - "--create-namespace"
`)
	// validDeployments = Deployments{
	// 	Helm: []HelmDeployment{
	// 		{
	// 			Option{
	// 				Arguments:   "jaeger",
	// 				Description: "Deploys Jaeger to your Kubernetes cluster",
	// 			},
	// 			HelmChart{
	// 				ChartName:   "jaeger",
	// 				ReleaseName: "jaeger",
	// 				Namespace:   "observability",
	// 				ChartPath:   "jaegertracing/jaeger",
	// 				Version:     "0.39.5",
	// 				Values:      []string{"/path/to/values/jaeger/jaeger-values.yaml"},
	// 				SetValues:   []string{"image=test:0.1.0"},
	// 				InstallArgs: []string{"--create-namespace"},
	// 				Repository: HelmRepository{
	// 					RepoName:    "jaegertracing",
	// 					RepoUrl:     "https://jaegertracing.github.io/helm-charts",
	// 					PrivateRepo: false,
	// 				},
	// 			},
	// 			1,
	// 		},
	// 	},
	// }
)

func importConfig(cfg []byte) {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer((cfg)))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func TestUnmarshalKruiseConfig(t *testing.T) {
	importConfig(validConfig)
	var kfg KruiseConfig
	err := viper.UnmarshalExact(&kfg)
	assert.NoError(t, err, "Unable to unmarshal KruiseConfig")
}
