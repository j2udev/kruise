package kruise

import (
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
	Files    []string
	Override string
	Manifest latest.KruiseConfig
}

// NewKonfig is used to create a new Kruise config (Konfig) object
//
// If the KRUISE_CONFIG environment variable is set, that config file is used,
// otherwise the following locations are checked in this order:
//
// cwd/kruise.yaml
//
// xdg.ConfigHome/kruise/kruise.yaml
//
// xdg.Home/.kruise.yaml
func NewKonfig() *Konfig {
	cfg := new(Konfig)
	cwd, err := os.Getwd()
	Fatal(err)
	cfg.Files = []string{
		cwd + "/kruise.yaml",
		xdg.ConfigHome + "/kruise/kruise.yaml",
		xdg.Home + "/.kruise.yaml",
	}
	cfg.Override = os.Getenv("KRUISE_CONFIG")
	cfg.ApplyUserConfig()
	return cfg
}

// ApplyUserConfig reads in a configuration file that is passed to viper and
// unmarshalled
func (k *Konfig) ApplyUserConfig() {
	Logger.Debug("Setting user config")
	k.setUserConfig()
	Logger.Debug("Unmarshalling user config")
	k.unmarshalExactConfig()
}

func (k Konfig) setUserConfig() {
	if k.Override != "" {
		Logger.Debugf("Overriding default user config: %s", k.Override)
		viper.SetConfigFile(k.Override)
		Logger.Debugf("Default user config overridden: %s", k.Override)
	} else {
		for _, file := range k.Files {
			Logger.Debugf("Searching for config in %s", file)
			if exists(file) {
				Logger.Debugf("Config found in %s", file)
				viper.SetConfigFile(file)
				break
			}
			Logger.Debugf("Config not found in %s", file)
		}
	}
	k.readConfig()
}

func (k Konfig) readConfig() {
	if err := viper.ReadInConfig(); err != nil {
		if k.Override != "" {
			Logger.Warnf("No user supplied config found in: %v", k.Override)
		} else {
			Logger.Warnf("No user supplied config found in the following paths: %v", k.Files)
		}
	}
}

func (k *Konfig) unmarshalExactConfig() {
	if err := viper.UnmarshalExact(&k.Manifest); err != nil {
		Fatalf(err, "Unable to decode config into struct")
	}
	logger := k.Manifest.Logger
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
	}
	Logger.Infof("Using config file: %s", viper.ConfigFileUsed())
	Logger.Debug("Config successfully unmarshalled!")
}
