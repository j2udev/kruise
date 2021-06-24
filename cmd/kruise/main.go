package main

import (
	"github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/spf13/cobra"
)

// Simple main function that creates a new kruise command
func main() {
	cobra.CheckErr(kruise.NewKruiseCmd().Execute())
}
