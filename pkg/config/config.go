package config

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigFile struct used to make passing around config files easier
type ConfigFile struct {
	Path      string
	Extension string
	FileName  string
	Override  string
}

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

var config Config

// InitConfig initializes default config
func InitConfig() {
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	configFile := ConfigFile{
		Path:      home,
		Extension: "yaml",
		FileName:  ".kruise",
	}
	InitCustomConfig(configFile)
}

// InitCustomConfig reads in a ConfigFile that is passed to viper
func InitCustomConfig(configFile ConfigFile) {
	if configFile.Override != "" {
		// Use config file from override
		viper.SetConfigFile(configFile.Override)
	} else {
		viper.AddConfigPath(configFile.Path)
		viper.SetConfigType(configFile.Extension)
		viper.SetConfigName(configFile.FileName)
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode config into struct, %v", err)
	}
}
