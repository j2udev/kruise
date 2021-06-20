package config

// Config struct used to unmarshal yaml kruise configuration
type Config struct {
	Deploy map[string][]DeployConfig `mapstructure:"deploy"`
}


// DeployConfig struct used to unmarshal nested yaml kruise configuration
type DeployConfig struct {
	Mongodb map[string][]HelmConfig `mapstructure:"mongodb"`
	Kafka   map[string][]HelmConfig `mapstructure:"kafka"`
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
