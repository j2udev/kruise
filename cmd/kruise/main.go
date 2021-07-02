package main

import (
	c "github.com/j2udevelopment/kruise/pkg/config"
	k "github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/spf13/cobra"
)

// Simple main function that creates a new kruise command
func main() {
	c.InitConfig()
	cobra.CheckErr(k.NewKruiseCmd().Execute())
}
