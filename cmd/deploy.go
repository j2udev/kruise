package cmd

import (
	c "github.com/j2udevelopment/kruise/pkg/config"
	"github.com/j2udevelopment/kruise/pkg/helm"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

var chartVersion string
var chartNamespace string

var jaegerDeployment = c.HelmConfig{ReleaseName: "jaeger", ChartPath: "jaegertracing/jaeger"}
var kafkaDeployment = c.HelmConfig{ReleaseName: "kafka", ChartPath: "bitnami/kafka"}
var mongoDeployment = c.HelmConfig{ReleaseName: "mongodb", ChartPath: "bitnami/mongodb"}
var mysqlDeployment = c.HelmConfig{ReleaseName: "mysql", ChartPath: "bitnami/mysql"}
var postgresqlDeployment = c.HelmConfig{ReleaseName: "postgresql", ChartPath: "bitnami/postgresql"}
var prometheusOperatorDeployment = c.HelmConfig{ReleaseName: "prometheus-operator", ChartPath: "prometheus-community/kube-prometheus-stack"}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:       "deploy",
	Short:     "Deploy the specified options to your Kubernetes cluster",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{"jaeger", "kafka", "mongodb", "mysql", "postgresql", "prometheus-operator"},
	Run: func(cmd *cobra.Command, args []string) {
		wg := sync.WaitGroup{}
		wg.Add(len(args))
		for i := 0; i < len(args); i++ {
			deployment := strings.ToLower(args[i])
			switch {
			case deployment == "jaeger":
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.jaeger"), &jaegerDeployment)
					helm.Install(jaegerDeployment)
					wg.Done()
				}()
			case deployment == "kafka":
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.kafka"), &kafkaDeployment)
					helm.Install(kafkaDeployment)
					wg.Done()
				}()
			case deployment == "mongodb":
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.mongodb"), &mongoDeployment)
					helm.Install(mongoDeployment)
					wg.Done()
				}()
			case deployment == "mysql":
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.mysql"), &mysqlDeployment)
					helm.Install(mysqlDeployment)
					wg.Done()
				}()
			case deployment == "postgresql":
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.postgresql"), &postgresqlDeployment)
					helm.Install(postgresqlDeployment)
					wg.Done()
				}()
			case deployment == "prometheus-operator":
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.prometheus-operator"), &prometheusOperatorDeployment)
					helm.Install(prometheusOperatorDeployment)
					wg.Done()
				}()
			}
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.PersistentFlags().StringVarP(&chartVersion, "version", "v", "", "Override the Helm chart version for the specified deployments")
	deployCmd.PersistentFlags().StringVarP(&chartNamespace, "namespace", "n", "", "Override the namespace for the specified deployments")
	deployCmd.PersistentFlags().StringVar(&jaegerDeployment.Version, "jaeger-version", "", "Override the Jaeger Helm chart version")
	deployCmd.PersistentFlags().StringVar(&jaegerDeployment.Namespace, "jaeger-namespace", "observability", "Override the Jaeger namespace")
	deployCmd.PersistentFlags().StringVar(&kafkaDeployment.Version, "kafka-version", "", "Override the Kafka Helm chart version")
	deployCmd.PersistentFlags().StringVar(&kafkaDeployment.Namespace, "kafka-namespace", "kafka", "Override the Kafka namespace")
	deployCmd.PersistentFlags().StringVar(&mongoDeployment.Version, "mongodb-version", "", "Override the MongoDB Helm chart version")
	deployCmd.PersistentFlags().StringVar(&mongoDeployment.Namespace, "mongodb-namespace", "mongodb", "Override the MongoDB namespace")
	deployCmd.PersistentFlags().StringVar(&mysqlDeployment.Version, "mysql-version", "", "Override the MySQL Helm chart version")
	deployCmd.PersistentFlags().StringVar(&mysqlDeployment.Namespace, "mysql-namespace", "mysql", "Override the MySQL namespace")
	deployCmd.PersistentFlags().StringVar(&postgresqlDeployment.Version, "postgresql-version", "", "Override the PostgreSQL Helm chart version")
	deployCmd.PersistentFlags().StringVar(&postgresqlDeployment.Namespace, "postgresql-namespace", "postgresql", "Override the PostgreSQL namespace")
	deployCmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Version, "prom-op-version", "", "Override the Prometheus Operator Helm chart version")
	deployCmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Namespace, "prom-op-namespace", "monitoring", "Override the Prometheus Operator namespace")
}
