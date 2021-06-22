package cmd

import (
	"sync"

	c "github.com/j2udevelopment/kruise/pkg/config"
	h "github.com/j2udevelopment/kruise/pkg/helm"
	u "github.com/j2udevelopment/kruise/pkg/utils"
	t "github.com/j2udevelopment/kruise/tpl"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deployOpts = []u.Option{
	{
		Arguments:   "jaeger",
		Description: "Deploys Jaeger to your Kubernetes Cluster",
	},
	{
		Arguments:   "kafka",
		Description: "Deploys Kafka to your Kubernetes Cluster",
	},
	{
		Arguments:   "mongodb, mongo",
		Description: "Deploys MongoDB to your Kubernetes Cluster",
	},
	{
		Arguments:   "mysql",
		Description: "Deploys MySQL to your Kubernetes Cluster",
	},
	{
		Arguments:   "postgresql, postgres",
		Description: "Deploys PostgreSQL to your Kubernetes Cluster",
	},
	{
		Arguments:   "prometheus-operator, prom-op",
		Description: "Deploys Prometheus Operator to your Kubernetes Cluster",
	},
}

var chartVersion string
var chartNamespace string
var jaegerDeployment = c.HelmConfig{ReleaseName: "jaeger", ChartPath: "jaegertracing/jaeger"}
var kafkaDeployment = c.HelmConfig{ReleaseName: "kafka", ChartPath: "bitnami/kafka"}
var mongodbDeployment = c.HelmConfig{ReleaseName: "mongodb", ChartPath: "bitnami/mongodb"}
var mysqlDeployment = c.HelmConfig{ReleaseName: "mysql", ChartPath: "bitnami/mysql"}
var postgresqlDeployment = c.HelmConfig{ReleaseName: "postgresql", ChartPath: "bitnami/postgresql"}
var prometheusOperatorDeployment = c.HelmConfig{ReleaseName: "prometheus-operator", ChartPath: "prometheus-community/kube-prometheus-stack"}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:       "deploy",
	Short:     "Deploy the specified options to your Kubernetes cluster",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: u.CollectValidArgs(deployOpts),
	Run: func(cmd *cobra.Command, args []string) {
		validArgsMap := u.CollectValidArgsMap(deployOpts)
		wg := sync.WaitGroup{}
		wg.Add(len(args))
		for _, arg := range args {
			switch {
			case u.Contains(validArgsMap["jaeger"], arg):
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.jaeger"), &jaegerDeployment)
					h.Install(jaegerDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["kafka"], arg):
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.kafka"), &kafkaDeployment)
					h.Install(kafkaDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["mongodb"], arg):
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.mongodb"), &mongodbDeployment)
					h.Install(mongodbDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["mysql"], arg):
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.mysql"), &mysqlDeployment)
					h.Install(mysqlDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["postgresql"], arg):
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.postgresql"), &postgresqlDeployment)
					h.Install(postgresqlDeployment)
					wg.Done()
				}()
			case u.Contains(validArgsMap["prometheus-operator"], arg):
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.prometheus-operator"), &prometheusOperatorDeployment)
					h.Install(prometheusOperatorDeployment)
					wg.Done()
				}()
			}
		}
		wg.Wait()
	},
}

func init() {
	var wrapper = u.CommandWrapper{
		Cmd:  deployCmd,
		Opts: deployOpts,
	}
	rootCmd.AddCommand(deployCmd)
	deployCmd.SetUsageTemplate(t.UsageTemplate())
	deployCmd.SetUsageFunc(t.UsageFunc(wrapper))
	deployCmd.PersistentFlags().StringVarP(&chartVersion, "version", "v", "", "Override the Helm chart version for the specified deployments")
	deployCmd.PersistentFlags().StringVarP(&chartNamespace, "namespace", "n", "", "Override the namespace for the specified deployments")
	deployCmd.PersistentFlags().StringVar(&jaegerDeployment.Version, "jaeger-version", "", "Override the Jaeger Helm chart version")
	deployCmd.PersistentFlags().StringVar(&jaegerDeployment.Namespace, "jaeger-namespace", "observability", "Override the Jaeger namespace")
	deployCmd.PersistentFlags().StringVar(&kafkaDeployment.Version, "kafka-version", "", "Override the Kafka Helm chart version")
	deployCmd.PersistentFlags().StringVar(&kafkaDeployment.Namespace, "kafka-namespace", "kafka", "Override the Kafka namespace")
	deployCmd.PersistentFlags().StringVar(&mongodbDeployment.Version, "mongodb-version", "", "Override the MongoDB Helm chart version")
	deployCmd.PersistentFlags().StringVar(&mongodbDeployment.Namespace, "mongodb-namespace", "mongodb", "Override the MongoDB namespace")
	deployCmd.PersistentFlags().StringVar(&mysqlDeployment.Version, "mysql-version", "", "Override the MySQL Helm chart version")
	deployCmd.PersistentFlags().StringVar(&mysqlDeployment.Namespace, "mysql-namespace", "mysql", "Override the MySQL namespace")
	deployCmd.PersistentFlags().StringVar(&postgresqlDeployment.Version, "postgresql-version", "", "Override the PostgreSQL Helm chart version")
	deployCmd.PersistentFlags().StringVar(&postgresqlDeployment.Namespace, "postgresql-namespace", "postgresql", "Override the PostgreSQL namespace")
	deployCmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Version, "prom-op-version", "", "Override the Prometheus Operator Helm chart version")
	deployCmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Namespace, "prom-op-namespace", "monitoring", "Override the Prometheus Operator namespace")
}
