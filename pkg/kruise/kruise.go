package kruise

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Kfg    *Konfig
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
