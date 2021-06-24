package deploy

import (
	"sync"

	"github.com/fatih/color"
	c "github.com/j2udevelopment/kruise/pkg/config"
	h "github.com/j2udevelopment/kruise/pkg/helm"
	u "github.com/j2udevelopment/kruise/pkg/utils"
	t "github.com/j2udevelopment/kruise/tpl"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chartVersion string
var chartNamespace string
var jaegerDeployment = c.HelmConfig{ReleaseName: "jaeger", ChartPath: "jaegertracing/jaeger"}
var kafkaDeployment = c.HelmConfig{ReleaseName: "kafka", ChartPath: "bitnami/kafka"}
var mongodbDeployment = c.HelmConfig{ReleaseName: "mongodb", ChartPath: "bitnami/mongodb"}
var mysqlDeployment = c.HelmConfig{ReleaseName: "mysql", ChartPath: "bitnami/mysql"}
var postgresqlDeployment = c.HelmConfig{ReleaseName: "postgresql", ChartPath: "bitnami/postgresql"}
var prometheusOperatorDeployment = c.HelmConfig{ReleaseName: "prometheus-operator", ChartPath: "prometheus-community/kube-prometheus-stack"}

// NewDeployOpts returns options for the deploy command
func NewDeployOpts() []u.Option {
	opts := []u.Option{
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
	return opts
}

// NewDeployCmd represents the deploy command
func NewDeployCmd(opts []u.Option) *cobra.Command {
	shallowDryRun := false
	cmd := &cobra.Command{
		Use:       "deploy",
		Short:     "Deploy the specified options to your Kubernetes cluster",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: u.CollectValidArgs(opts),
		Run: func(cmd *cobra.Command, args []string) {
			validArgsMap := u.CollectValidArgsMap(opts)
			wg := sync.WaitGroup{}
			wg.Add(len(args))
			for _, arg := range args {
				switch {
				case u.Contains(validArgsMap["jaeger"], arg):
					go func() {
						mapstructure.Decode(viper.GetStringMap("deploy.jaeger"), &jaegerDeployment)
						color.Green("Attempting to deploy Jaeger...")
						if err := h.Install(shallowDryRun, &jaegerDeployment); err != nil {
							color.Red("Failed to deploy Jaeger")
						}
						wg.Done()
					}()
				case u.Contains(validArgsMap["kafka"], arg):
					go func() {
						mapstructure.Decode(viper.GetStringMap("deploy.kafka"), &kafkaDeployment)
						color.Green("Attempting to deploy Kafka...")
						if err := h.Install(shallowDryRun, &kafkaDeployment); err != nil {
							color.Red("Failed to deploy Kafka")
						}
						wg.Done()
					}()
				case u.Contains(validArgsMap["mongodb"], arg):
					go func() {
						mapstructure.Decode(viper.GetStringMap("deploy.mongodb"), &mongodbDeployment)
						color.Green("Attempting to deploy MongoDB...")
						if err := h.Install(shallowDryRun, &mongodbDeployment); err != nil {
							color.Red("Failed to deploy MongoDB")
						}
						wg.Done()
					}()
				case u.Contains(validArgsMap["mysql"], arg):
					go func() {
						mapstructure.Decode(viper.GetStringMap("deploy.mysql"), &mysqlDeployment)
						color.Green("Attempting to deploy MySQL...")
						if err := h.Install(shallowDryRun, &mysqlDeployment); err != nil {
							color.Red("Failed to deploy MySQL")
						}
						wg.Done()
					}()
				case u.Contains(validArgsMap["postgresql"], arg):
					go func() {
						mapstructure.Decode(viper.GetStringMap("deploy.postgresql"), &postgresqlDeployment)
						color.Green("Attempting to deploy PostgreSQL...")
						if err := h.Install(shallowDryRun, &postgresqlDeployment); err != nil {
							color.Red("Failed to deploy PostgreSQL")
						}
						wg.Done()
					}()
				case u.Contains(validArgsMap["prometheus-operator"], arg):
					go func() {
						mapstructure.Decode(viper.GetStringMap("deploy.prometheus-operator"), &prometheusOperatorDeployment)
						color.Green("Attempting to deploy Prometheus Operator...")
						if err := h.Install(shallowDryRun, &prometheusOperatorDeployment); err != nil {
							color.Red("Failed to deploy Prometheus Operator")
						}
						wg.Done()
					}()
				}
			}
			wg.Wait()
		},
	}
	wrapper := u.CommandWrapper{
		Cmd:  cmd,
		Opts: opts,
	}
	cmd.SetUsageTemplate(t.UsageTemplate())
	cmd.SetUsageFunc(t.UsageFunc(wrapper))
	cmd.PersistentFlags().StringVarP(&chartVersion, "version", "v", "", "Override the Helm chart version for the specified deployments")
	cmd.PersistentFlags().StringVarP(&chartNamespace, "namespace", "n", "", "Override the namespace for the specified deployments")
	cmd.PersistentFlags().StringVar(&jaegerDeployment.Version, "jaeger-version", "", "Override the Jaeger Helm chart version")
	cmd.PersistentFlags().StringVar(&jaegerDeployment.Namespace, "jaeger-namespace", "observability", "Override the Jaeger namespace")
	cmd.PersistentFlags().StringVar(&kafkaDeployment.Version, "kafka-version", "", "Override the Kafka Helm chart version")
	cmd.PersistentFlags().StringVar(&kafkaDeployment.Namespace, "kafka-namespace", "kafka", "Override the Kafka namespace")
	cmd.PersistentFlags().StringVar(&mongodbDeployment.Version, "mongodb-version", "", "Override the MongoDB Helm chart version")
	cmd.PersistentFlags().StringVar(&mongodbDeployment.Namespace, "mongodb-namespace", "mongodb", "Override the MongoDB namespace")
	cmd.PersistentFlags().StringVar(&mysqlDeployment.Version, "mysql-version", "", "Override the MySQL Helm chart version")
	cmd.PersistentFlags().StringVar(&mysqlDeployment.Namespace, "mysql-namespace", "mysql", "Override the MySQL namespace")
	cmd.PersistentFlags().StringVar(&postgresqlDeployment.Version, "postgresql-version", "", "Override the PostgreSQL Helm chart version")
	cmd.PersistentFlags().StringVar(&postgresqlDeployment.Namespace, "postgresql-namespace", "postgresql", "Override the PostgreSQL namespace")
	cmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Version, "prom-op-version", "", "Override the Prometheus Operator Helm chart version")
	cmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Namespace, "prom-op-namespace", "monitoring", "Override the Prometheus Operator namespace")

	return cmd
}
