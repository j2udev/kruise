package cmd

import (
	"strings"
	"sync"

	h "github.com/j2udevelopment/kruise/pkg/helm"
	u "github.com/j2udevelopment/kruise/pkg/utils"
	"github.com/spf13/cobra"
)

var deleteOpts = []u.Option{
	{
		Arguments:   "jaeger",
		Description: "Deletes Jaeger from your Kubernetes Cluster",
	},
	{
		Arguments:   "kafka",
		Description: "Deletes Kafka from your Kubernetes Cluster",
	},
	{
		Arguments:   "mongodb, mongo",
		Description: "Deletes MongoDB from your Kubernetes Cluster",
	},
	{
		Arguments:   "mysql",
		Description: "Deletes MySQL from your Kubernetes Cluster",
	},
	{
		Arguments:   "postgresql, postgres",
		Description: "Deletes PostgreSQL from your Kubernetes Cluster",
	},
	{
		Arguments:   "prometheus-operator, prom-op",
		Description: "Deletes Prometheus Operator from your Kubernetes Cluster",
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:       "delete",
	Short:     "Delete the specified options from your Kubernetes cluster",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: u.CollectValidArgs(deleteOpts),
	Run: func(cmd *cobra.Command, args []string) {
		validArgsMap := u.CollectValidArgsMap(deleteOpts)
		wg := sync.WaitGroup{}
		wg.Add(len(args))
		for _, arg := range args {
			switch {
			case u.Contains(validArgsMap["jaeger"], arg):
				go func() {
					h.Uninstall(jaegerDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["kafka"], arg):
				go func() {
					h.Uninstall(kafkaDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["mongodb"], arg):
				go func() {
					h.Uninstall(mongodbDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["mysql"], arg):
				go func() {
					h.Uninstall(mysqlDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["postgresql"], arg):
				go func() {
					h.Uninstall(postgresqlDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["prometheus-operator"], arg):
				go func() {
					h.Uninstall(prometheusOperatorDeployment)
					wg.Done()
				}()
			}
		}

		for i := 0; i < len(args); i++ {
			deployment := strings.ToLower(args[i])
			switch {
			case deployment == "jaeger":
				go func() {
					h.Uninstall(jaegerDeployment)
					wg.Done()
				}()
			case deployment == "kafka":
				go func() {
					h.Uninstall(kafkaDeployment)
					wg.Done()
				}()
			case deployment == "mongodb":
				go func() {
					h.Uninstall(mongodbDeployment)
					wg.Done()
				}()
			case deployment == "mysql":
				go func() {
					h.Uninstall(mysqlDeployment)
					wg.Done()
				}()
			case deployment == "postgresql":
				go func() {
					h.Uninstall(postgresqlDeployment)
					wg.Done()
				}()
			case deployment == "prometheus-operator":
				go func() {
					h.Uninstall(prometheusOperatorDeployment)
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
	deleteCmd.PersistentFlags().StringVar(&mongodbDeployment.Namespace, "mongodb-namespace", "mongodb", "Override the MongoDB namespace")
	deleteCmd.PersistentFlags().StringVar(&mysqlDeployment.Namespace, "mysql-namespace", "mysql", "Override the MySQL namespace")
	deleteCmd.PersistentFlags().StringVar(&postgresqlDeployment.Namespace, "postgresql-namespace", "postgresql", "Override the PostgreSQL namespace")
	deleteCmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Namespace, "prom-op-namespace", "monitoring", "Override the Prometheus Operator namespace")
}
