package config

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// File struct used to capture config file information to be passed to viper
type File struct {
	Path      string
	Extension string
	FileName  string
	Override  string
}

// Manifest struct used to unmarshal yaml kruise configuration
type Manifest struct {
	Deploy map[string][]Deployments `mapstructure:"deploy"`
	Delete map[string][]Deployments `mapstructure:"delete"`
}

// Deployments struct used to unmarshal yaml kruise configuration
type Deployments struct {
	Helm map[string][]HelmDeployment `mapstructure:"helm"`
}

// HelmDeployment struct used to unmarshal yaml kruise configuration
type HelmDeployment struct {
	Option      `mapstructure:"option"`
	HelmCommand `mapstructure:"command"`
}

// HelmCommand struct used to unmarshal yaml kruise configuration
type HelmCommand struct {
	ReleaseName string
	ChartPath   string
	Namespace   string
	Version     string
	Values      []string
	Args        []string
	ExtraArgs   []string
}

// Option struct used to unmarshal yaml kruise configuration and facilitate
// wrapping cobra commands
type Option struct {
	Arguments   string
	Description string
}

// CommandWrapper is used to wrap cobra commands to support command options
type CommandWrapper struct {
	Cmd  *cobra.Command
	Opts *[]Option
}

// Initialize reads in a configuration file that is passed to viper and
// unmarshalled
func Initialize(configFile File, data interface{}) {
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
	} else {
		log.Fatalln("Something is wrong with the config path:", err)
	}
	err := viper.Unmarshal(&data)
	if err != nil {
		log.Fatalf("Unable to decode config into struct, %v", err)
	}
}

// Decode is used to destructure config maps into structs
func Decode(key string, data interface{}) error {
	return mapstructure.Decode(viper.Get(key), data)
}
