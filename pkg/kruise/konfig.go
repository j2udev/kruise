package kruise

import (
	"fmt"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/spf13/viper"
)

type (
	// Konfig struct used to combine file metadata with unmarshalled Kruise
	// configuration
	Konfig struct {
		Metadata Metadata
		Manifest latest.KruiseConfig
	}

	// Metadata struct used to capture config file information to be passed to viper
	Metadata struct {
		Paths     []string
		Extension string
		Name      string
		Override  string
	}
)

func NewKonfig() *Konfig {
	cfg := new(Konfig)
	meta := Metadata{
		Paths: []string{
			xdg.ConfigHome + "/kruise",
			xdg.Home,
		},
		Name:      ".kruise",
		Extension: "yaml",
	}
	cfg.Metadata = meta
	return cfg
}

// Initialize reads in a configuration file that is passed to viper and
// unmarshalled
func (kfg *Konfig) Initialize() {
	if kfg.Metadata.Override != "" {
		// Use config file from override
		viper.SetConfigFile(kfg.Metadata.Override)
	} else {
		viper.SetConfigName(kfg.Metadata.Name)
		viper.SetConfigType(kfg.Metadata.Extension)
		for _, path := range kfg.Metadata.Paths {
			viper.AddConfigPath(path)
		}
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatalln("Something is wrong with the config path:", err)
	}
	err := viper.UnmarshalExact(&kfg.Manifest)
	if err != nil {
		log.Fatalf("Unable to decode config into struct, %v", err)
	}
}

// func (kfg Konfig) GetDeployConfig() Deployments {
// 	return Deployments{kfg.Manifest.Deploy}
// }

// func (kfg Konfig) GetDeleteConfig() Deployments {
// 	return Deployments{kfg.Manifest.Delete}
// }
