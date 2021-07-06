package main

import (
	"github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/spf13/cobra"
)

// Simple main function that creates a new kruise command
func main() {
	cobra.OnInitialize(kruise.Initialize)
	cobra.CheckErr(kruise.NewKruiseCmd().Execute())

	//TODO: Figure out why the help flag is not initializing kruise config,
	// options, etc
	//TODO: The help command and usage works perfectly fine
	//TODO: Abstract this kind of testing to a dedicated test file

	// home, _ := homedir.Dir()

	// //non-working args
	// args := []string{"deploy", "-h", fmt.Sprintf("--config=%s", home+"/.kruise.yaml")}

	// //working args
	// args := []string{"help", "deploy", fmt.Sprintf("--config=%s", home+"/.kruise-2.yaml")}

	// cmd := kruise.NewKruiseCmd()
	// cmd.SetArgs(args)
	// cobra.CheckErr(cmd.Execute())
}
