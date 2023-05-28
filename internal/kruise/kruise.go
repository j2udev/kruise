// Kruise is a black box CLI that is driven by a config
//
// Kruise helps with abstractly deploying workloads to Kubernetes
package kruise

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// Kfg is a global config object for Kruise into which config is unmarshalled
	Kfg *Konfig
	// Logger is the global Logrus logger used by Kruise
	//
	// Read more about Logrus: https://pkg.go.dev/github.com/sirupsen/logrus
	Logger *logrus.Logger
)

// Initialize is used to initialize Kruise
func Initialize() {
	InitializeLogger()
	InitializeConfig()
	Logger.Trace("Kruise initialized")
}

// InitializeConfig is used to initialize Kruise configuration
func InitializeConfig() {
	Kfg = NewKonfig()
	Logger.Trace("Config initialized")
}

// InitializeLogger is used to initialize the Kruise logger
func InitializeLogger() {
	logger := logrus.New()
	logger.Out = os.Stdout
	Logger = logger
	Logger.Trace("Logger initialized")
}
