package kruise

import (
	"os"
)

var Kfg Konfig

// Initialize is used to initialize kruise configuration and command options
func Initialize() {
	home, err := os.UserHomeDir()
	CheckErr(err)
	file := &Kfg.Metadata
	file.Paths = []string{home, home + ".config/kruise"}
	file.Name = ".kruise"
	file.Extension = "yaml"
	Kfg.Initialize()
}
