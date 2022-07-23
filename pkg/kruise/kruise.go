package kruise

import (
	"github.com/sirupsen/logrus"
)

var Kfg *Konfig
var Logger *logrus.Logger

// Initialize is used to initialize Kruise
func Initialize() {
	InitializeLogger()
	InitializeConfig()
}

// InitializeConfig is used to initialize Kruise configuration
func InitializeConfig() {
	Kfg = NewKonfig()
	Kfg.Initialize()
}

// InitializeLogger is used to initialize the Kruise logger
func InitializeLogger() {
	logger := logrus.New()
	Logger = logger
}
