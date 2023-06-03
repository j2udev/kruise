// Kruise is a black box CLI that is driven by a config
//
// Kruise helps with abstractly deploying workloads to Kubernetes
package kruise

import (
	"os"

	"github.com/charmbracelet/log"
)

var (
	// Kfg is a global config object for Kruise into which config is unmarshalled
	Kfg *Konfig
	// Logger is the global logger used by Kruise
	Logger *log.Logger
)

// Initialize is used to initialize Kruise
func Initialize() {
	InitializeLogger()
	InitializeConfig()
	Logger.Debug("Kruise initialized")
}

// InitializeConfig is used to initialize Kruise configuration
func InitializeConfig() {
	Kfg = NewKonfig()
	Logger.Debug("Config initialized")
}

// InitializeLogger is used to initialize the Kruise logger
func InitializeLogger() {
	Logger = log.New(os.Stderr)
	Logger.SetLevel(log.WarnLevel)
}
