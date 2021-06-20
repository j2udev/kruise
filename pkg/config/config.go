package config

type Config struct {
	Deploy map[string][]DeployConfig `mapstructure:"deploy"`
}

type DeployConfig struct {
	Mongodb map[string][]HelmConfig `mapstructure:"mongodb"`
	Kafka   map[string][]HelmConfig `mapstructure:"kafka"`
}

type HelmConfig struct {
	ReleaseName string
	ChartPath   string
	Namespace   string
	Version     string
	Values      []string
	Args        []string
	ExtraArgs   []string
}
