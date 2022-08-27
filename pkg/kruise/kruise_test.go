package kruise

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
)

type ObservabilityIntTestSuite struct {
	suite.Suite
	kfg *Konfig
	fs  *pflag.FlagSet
}

func TestObservabilityIntTestSuite(t *testing.T) {
	suite.Run(t, &ObservabilityIntTestSuite{})
}

func (s *ObservabilityIntTestSuite) SetupSuite() {
	s.T().Log("Setting Up Observability Integration Test Suite")
	os.Setenv("KRUISE_CONFIG", "../../examples/observability/kruise.yaml")
	Initialize()
	s.kfg = Kfg
	s.fs = GetDeployFlags()
	err := s.fs.Set("shallow-dry-run", "true")
	if err != nil {
		s.FailNow(err.Error())
	}
}

func (s *ObservabilityIntTestSuite) TestIstioDeployment() {
	actual := trimDeployStdoutPrefix(s.deployIstio)
	s.Equal(s.expectedIstio(), actual)
}

func (s *ObservabilityIntTestSuite) TestJaegerDeployment() {
	actual := trimDeployStdoutPrefix(s.deployJaeger)
	s.Equal(s.expectedJaeger(), actual)
}

func (s *ObservabilityIntTestSuite) TestLokiDeployment() {
	actual := trimDeployStdoutPrefix(s.deployLoki)
	s.Equal(s.expectedLoki(), actual)
}

func (s *ObservabilityIntTestSuite) TestPrometheusOperatorDeployment() {
	actual := trimDeployStdoutPrefix(s.deployPrometheusOperator)
	s.Equal(s.expectedPrometheusOperator(), actual)
}

func (s *ObservabilityIntTestSuite) TestObservabilityProfile() {
	actual := trimDeployStdoutPrefix(s.deployObservability)
	s.Equal(s.expectedObservability(), actual)
	actual = trimDeployStdoutPrefix(s.deployTelemetry)
	s.Equal(s.expectedObservability(), actual)
}

func (s *ObservabilityIntTestSuite) deployIstio() {
	Deploy(s.fs, []string{"istio"})
}

func (s *ObservabilityIntTestSuite) deployJaeger() {
	Deploy(s.fs, []string{"jaeger"})
}

func (s *ObservabilityIntTestSuite) deployLoki() {
	Deploy(s.fs, []string{"loki"})
}

func (s *ObservabilityIntTestSuite) deployPrometheusOperator() {
	Deploy(s.fs, []string{"prometheus-operator"})
}

func (s *ObservabilityIntTestSuite) deployObservability() {
	Deploy(s.fs, []string{"observability"})
}

func (s *ObservabilityIntTestSuite) deployTelemetry() {
	Deploy(s.fs, []string{"telemetry"})
}

func (s *ObservabilityIntTestSuite) expectedIstio() string {
	expected := `
helm upgrade --install istio-base istio/base --namespace istio-system --version 1.14.1 -f values/istio-base-values.yaml --create-namespace
helm upgrade --install istiod istio/istiod --namespace istio-system --version 1.14.1 -f values/istiod-values.yaml --create-namespace
helm upgrade --install istio-ingressgateway istio/gateway --namespace istio-system --version 1.14.1 -f values/istio-gateway-values.yaml --set service.externalIPs[0]=CHANGE_ME --create-namespace
kubectl create namespace istio-system
kubectl apply --namespace istio-system -f manifests/istio-gateway.yaml
`
	expected = strings.TrimPrefix(expected, "\n")
	return expected
}

func (s *ObservabilityIntTestSuite) expectedJaeger() string {
	expected := `
helm upgrade --install jaeger jaegertracing/jaeger --namespace tracing --version 0.57.1 -f values/jaeger-values.yaml --create-namespace
kubectl create namespace tracing
kubectl apply --namespace tracing -f manifests/jaeger-virtual-service.yaml
`
	expected = strings.TrimPrefix(expected, "\n")
	return expected
}

func (s *ObservabilityIntTestSuite) expectedLoki() string {
	expected := `
helm upgrade --install loki grafana/loki-stack --namespace logging --version 2.6.5 --create-namespace
`
	expected = strings.TrimPrefix(expected, "\n")
	return expected
}

func (s *ObservabilityIntTestSuite) expectedPrometheusOperator() string {
	expected := `
helm upgrade --install prometheus-operator prometheus-community/kube-prometheus-stack --namespace monitoring --version 36.0.2 -f values/prometheus-operator-values.yaml --create-namespace
kubectl create namespace monitoring
kubectl apply --namespace monitoring -f manifests/grafana-virtual-service.yaml
`
	expected = strings.TrimPrefix(expected, "\n")
	return expected
}

func (s *ObservabilityIntTestSuite) expectedObservability() string {
	expected := s.expectedIstio() + s.expectedJaeger() + s.expectedLoki() + s.expectedPrometheusOperator()
	expected = strings.TrimPrefix(expected, "\n")
	return expected
}
