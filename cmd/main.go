package main

import (
	"github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/spf13/cobra"
)

// Simple main function that creates a new kruise command
func main() {
	cobra.OnInitialize(kruise.Initialize)
	cobra.CheckErr(kruise.NewCmd().Execute())

	//TODO: Figure out why the help flag is not initializing kruise config,
	// options, etc
	//TODO: The help command and usage works perfectly fine
	//TODO: Abstract this kind of testing to a dedicated test file

	// home, _ := homedir.Dir()
	// //non-working args
	// // notWorking := []string{"deploy", "--help", fmt.Sprintf("--config=%s", home+"/.kruise-2.yaml")}
	// // notWorking := []string{"deploy", "kafka", "--help", fmt.Sprintf("--config=%s", home+"/.kruise-2.yaml")}
	// // notWorking := []string{"deploy", fmt.Sprintf("--config=%s", home+"/.kruise-2.yaml")}

	// //working args
	// working := []string{"help", "deploy", fmt.Sprintf("--config=%s", home+"/.kruise-2.yaml")}
	// // working := []string{"deploy", "kafka", fmt.Sprintf("--config=%s", home+"/.kruise-2.yaml")}
	// // working := []string{"deploy", "kafka"}

	// // cmd := kruise.NewCmd()
	// // cmd.SetArgs(notWorking)
	// // // cmd.DisableFlagParsing = true
	// // cobra.CheckErr(cmd.Execute())
	// // // cmd.DebugFlags()
	// // // cmd.TraverseChildren = false

	// cmd := kruise.NewCmd()
	// cmd.SetArgs(working)
	// // cmd.DisableFlagParsing = true
	// cobra.CheckErr(cmd.Execute())
	// // cmd.DebugFlags()
}
