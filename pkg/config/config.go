package config

// Config struct used to unmarshal yaml kruise configuration
type Config struct {
	Deploy map[string][]DeployConfig `mapstructure:"deploy"`
}

// DeployConfig struct used to unmarshal nested yaml kruise configuration
type DeployConfig struct {
	Jaeger             map[string][]HelmConfig `mapstructure:"jaeger"`
	Kafka              map[string][]HelmConfig `mapstructure:"kafka"`
	Mongodb            map[string][]HelmConfig `mapstructure:"mongodb"`
	Mysql              map[string][]HelmConfig `mapstructure:"mysql"`
	Postgresql         map[string][]HelmConfig `mapstructure:"postgresql"`
	PrometheusOperator map[string][]HelmConfig `mapstructure:"prometheus-operator"`
}

// HelmConfig struct used to unmarshal nested yaml kruise configuration
type HelmConfig struct {
	ReleaseName string
	ChartPath   string
	Namespace   string
	Version     string
	Values      []string
	Args        []string
	ExtraArgs   []string
}
