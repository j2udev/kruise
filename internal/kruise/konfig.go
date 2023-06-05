package kruise

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
	"github.com/j2udev/kruise/internal/schema/latest"
	"github.com/spf13/viper"
)

// Konfig struct used to combine file metadata with unmarshalled Kruise
// configuration
type Konfig struct {
	Paths       []string
	HiddenPaths []string
	Name        string
	Override    string
	Manifest    latest.KruiseConfig
}

// NewKonfig is used to create a new Kruise config (Konfig) object
//
// If the KRUISE_CONFIG environment variable is set, that config file is used,
// otherwise the following locations are checked in this order:
//
// cwd/kruise.json/toml/yaml
//
// xdg.ConfigHome/kruise/kruise.json/toml/yaml
//
// cwd/.kruise.json/toml/yaml
//
// xdg.ConfigHome/kruise/.kruise.json/toml/yaml
//
// xdg.Home/.kruise.json/toml/yaml
func NewKonfig() *Konfig {
	cfg := new(Konfig)
	cwd, err := os.Getwd()
	if err != nil {
		Logger.Fatal(err)
	}
	cfg.Name = "kruise"
	cfg.Paths = []string{
		cwd,
		xdg.ConfigHome + "/kruise",
	}
	cfg.HiddenPaths = []string{
		cwd,
		xdg.ConfigHome + "/kruise",
		xdg.Home,
	}
	// The best we can do for overriding config file location is to allow setting
	// it with an environment variable
	// The CLI is driven by config so we can't override the config from the CLI
	cfg.Override = os.Getenv("KRUISE_CONFIG")
	cfg.ApplyUserConfig()
	return cfg
}

// ApplyUserConfig reads in a configuration file that is passed to viper and
// unmarshalled
func (k *Konfig) ApplyUserConfig() {
	Logger.Debug("Setting config")
	k.setConfig()
	Logger.Debug("Unmarshalling config")
	k.unmarshalConfig()
	if k.Override != "" {
		Logger.Infof("Using config file: %s", k.Override)
		return
	}
	Logger.Infof("Using config file: %s", viper.ConfigFileUsed())
}

// setConfig is used to set the kruise config file
func (k Konfig) setConfig() {
	if k.Override != "" {
		overrideConfig(k.Override)
		return
	}
	checkConfig(k)
}

// checkConfig is used to check the default locations for a config file and read
// it in if found
func checkConfig(k Konfig) {
	viper.SetConfigName(k.Name)
	for _, path := range k.Paths {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigName("." + k.Name)
		for _, path := range k.HiddenPaths {
			viper.AddConfigPath(path)
		}
		if err := viper.ReadInConfig(); err != nil {
			Logger.Warn(err)
		}
	}
}

// overrideConfig is used to read in viper config from a non-default
// location or URL
func overrideConfig(config string) {
	Logger.Debugf("Attempting to use config defined by KRUISE_CONFIG: %s", config)
	if strings.HasPrefix(config, "http") {
		// if config is set from a URL currently only yaml is supported; not sure if
		// there is a way around this
		// in order to read config from a URL, you must set the config type before
		// reading it in
		viper.SetConfigType("yaml")
		cfg, err := fetchConfigFromURL(config)
		if err != nil {
			Logger.Fatalf("Can't fetch config from %s", string(config))
		}
		if err := viper.ReadConfig(bytes.NewBuffer(cfg)); err != nil {
			Logger.Fatal(err)
		}
	} else {
		viper.SetConfigFile(config)
		if err := viper.ReadInConfig(); err != nil {
			Logger.Fatal(err)
		}
	}
}

// fetchConfigFromURL retrieves the configuration content from a remote URL
func fetchConfigFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to fetch config file: %s", resp.Status)
	}
	cfgContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return cfgContent, nil
}

// unmarshalConfig is used to unmarshal user defined config into Kruise
// schema
func (k *Konfig) unmarshalConfig() {
	if err := viper.UnmarshalExact(&k.Manifest); err != nil {
		Logger.Fatal(err)
	}
	logger := k.Manifest.Logger
	if logger.Caller {
		Logger.SetReportCaller(true)
	}
	if logger.TimeStamp {
		Logger.SetReportTimestamp(true)
	}
	if logger.TimeFormat != "" {
		// https://pkg.go.dev/time#pkg-constants
		switch strings.ToLower(logger.TimeFormat) {
		case "layout":
			Logger.SetTimeFormat(time.Layout)
		case "ansic":
			Logger.SetTimeFormat(time.ANSIC)
		case "unixdate":
			Logger.SetTimeFormat(time.UnixDate)
		case "rubydate":
			Logger.SetTimeFormat(time.RubyDate)
		case "rfc822":
			Logger.SetTimeFormat(time.RFC822)
		case "rfc822z":
			Logger.SetTimeFormat(time.RFC822Z)
		case "rfc850":
			Logger.SetTimeFormat(time.RFC850)
		case "rfc1123":
			Logger.SetTimeFormat(time.RFC1123)
		case "rfc1123z":
			Logger.SetTimeFormat(time.RFC1123Z)
		case "rfc3339":
			Logger.SetTimeFormat(time.RFC3339)
		case "rfc3339nano":
			Logger.SetTimeFormat(time.RFC3339Nano)
		case "kitchen":
			Logger.SetTimeFormat(time.Kitchen)
		case "stamp":
			Logger.SetTimeFormat(time.Stamp)
		case "stampmilli":
			Logger.SetTimeFormat(time.StampMilli)
		case "stampmicro":
			Logger.SetTimeFormat(time.StampMicro)
		case "stampnano":
			Logger.SetTimeFormat(time.StampNano)
		case "datetime":
			Logger.SetTimeFormat(time.DateTime)
		case "dateonly":
			Logger.SetTimeFormat(time.DateOnly)
		case "timeonly":
			Logger.SetTimeFormat(time.TimeOnly)
		default:
			Logger.Fatalf("Invalid time format %s\nValid time formats can be found here: https://pkg.go.dev/time#pkg-constants", logger.TimeFormat)
		}
		lvl := logger.Level
		if lvl != "" {
			switch lvl {
			case "debug":
				Logger.SetLevel(log.DebugLevel)
			case "info":
				Logger.SetLevel(log.InfoLevel)
			case "warn":
				Logger.SetLevel(log.WarnLevel)
			case "error":
				Logger.SetLevel(log.ErrorLevel)
			default:
				Logger.Fatalf("Invalid verbosity level: %s", lvl)
			}
		}
	}
	Logger.Debug("Config successfully unmarshalled!")
}
