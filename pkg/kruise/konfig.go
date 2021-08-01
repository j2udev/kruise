package kruise

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Konfig struct used to combine file metadata with unmarshalled kruise
// configuration
type Konfig struct {
	Metadata Metadata
	Manifest Manifest
}

// Metadata struct used to capture config file information to be passed to viper
type Metadata struct {
	Path      string
	Extension string
	Name      string
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

// Initialize reads in a configuration file that is passed to viper and
// unmarshalled
func (kfg Konfig) Initialize() {
	if kfg.Metadata.Override != "" {
		// Use config file from override
		viper.SetConfigFile(kfg.Metadata.Override)
	} else {
		viper.AddConfigPath(kfg.Metadata.Path)
		viper.SetConfigType(kfg.Metadata.Extension)
		viper.SetConfigName(kfg.Metadata.Name)
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatalln("Something is wrong with the config path:", err)
	}
	err := viper.Unmarshal(&kfg.Manifest)
	if err != nil {
		log.Fatalf("Unable to decode config into struct, %v", err)
	}
}

// Decode is used to destructure config maps into structs
func Decode(key string, data interface{}) error {
	return mapstructure.Decode(viper.Get(key), data)
}
