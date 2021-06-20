package cmd

import (
	"github.com/j2udevelopment/kruise/pkg/helm"
	"github.com/spf13/cobra"
	"strings"
	"sync"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:       "delete",
	Short:     "Delete the specified options from your Kubernetes cluster",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{"mongodb", "kafka"},
	Run: func(cmd *cobra.Command, args []string) {
		wg := sync.WaitGroup{}
		wg.Add(len(args))
		for i := 0; i < len(args); i++ {
			deployment := strings.ToLower(args[i])
			switch {
			case deployment == "mongodb":
				go func() {
					helm.Uninstall(mongoDeployment)
					wg.Done()
				}()
			case deployment == "kafka":
				go func() {
					helm.Uninstall(kafkaDeployment)
					wg.Done()
				}()
			}
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().StringVarP(&chartNamespace, "namespace", "n", "", "Override the namespace for the specified deployments")
	deleteCmd.PersistentFlags().StringVar(&kafkaDeployment.Namespace, "kafka-namespace", "kafka", "Override the Kafka namespace")
	deleteCmd.PersistentFlags().StringVar(&mongoDeployment.Namespace, "mongodb-namespace", "mongodb", "Override the MongoDB namespace")
}
