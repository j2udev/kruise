package main

import (
	"github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/spf13/cobra"
)

// Simple main function that creates a new kruise command
func main() {
	cobra.OnInitialize(kruise.Initialize)
	cobra.CheckErr(kruise.NewCmd().Execute())

	// -------------------------------------------------
	// ---------------------TESTING---------------------
	// -------------------------------------------------

	//TODO: Figure out why the help flag is not initializing kruise config
	//TODO: The help command and usage work perfectly fine
	//TODO: Abstract this kind of testing to a test file

	// var cmd *cobra.Command
	// home, _ := os.UserHomeDir()
	// //non-working args
	// notWorking := []string{"deploy", "--help", fmt.Sprintf("--config=%s", home+"/.kruise-2.yaml")}
	// cmd = kruise.NewCmd()
	// cmd.SetArgs(notWorking)
	// cobra.CheckErr(cmd.Execute())

	// //working args
	// working := []string{"help", "deploy", fmt.Sprintf("--config=%s", home+"/.kruise-2.yaml")}
	// cmd = kruise.NewCmd()
	// cmd.SetArgs(working)
	// cobra.CheckErr(cmd.Execute())
}
