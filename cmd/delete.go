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
			case deployment == "jaeger":
				go func() {
					helm.Uninstall(jaegerDeployment)
					wg.Done()
				}()
			case deployment == "kafka":
				go func() {
					helm.Uninstall(kafkaDeployment)
					wg.Done()
				}()
			case deployment == "mongodb":
				go func() {
					helm.Uninstall(mongoDeployment)
					wg.Done()
				}()
			case deployment == "mysql":
				go func() {
					helm.Uninstall(mysqlDeployment)
					wg.Done()
				}()
			case deployment == "postgresql":
				go func() {
					helm.Uninstall(postgresqlDeployment)
					wg.Done()
				}()
			case deployment == "prometheus-operator":
				go func() {
					helm.Uninstall(prometheusOperatorDeployment)
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
	deleteCmd.PersistentFlags().StringVar(&jaegerDeployment.Namespace, "jaeger-namespace", "observability", "Override the Jaeger namespace")
	deleteCmd.PersistentFlags().StringVar(&kafkaDeployment.Namespace, "kafka-namespace", "kafka", "Override the Kafka namespace")
	deleteCmd.PersistentFlags().StringVar(&mongoDeployment.Namespace, "mongodb-namespace", "mongodb", "Override the MongoDB namespace")
	deleteCmd.PersistentFlags().StringVar(&mysqlDeployment.Namespace, "mysql-namespace", "mysql", "Override the MySQL namespace")
	deleteCmd.PersistentFlags().StringVar(&postgresqlDeployment.Namespace, "postgresql-namespace", "postgresql", "Override the PostgreSQL namespace")
	deleteCmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Namespace, "prom-op-namespace", "monitoring", "Override the Prometheus Operator namespace")
}
