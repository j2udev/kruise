package cmd

import (
	c "github.com/j2udevelopment/kruise/pkg/config"
	"github.com/j2udevelopment/kruise/pkg/helm"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chartVersion string
var chartNamespace string

var jaegerDeployment = c.HelmConfig{ReleaseName: "jaeger", ChartPath: "jaegertracing/jaeger"}
var kafkaDeployment = c.HelmConfig{ReleaseName: "kafka", ChartPath: "bitnami/kafka"}
var mongoDeployment = c.HelmConfig{ReleaseName: "mongodb", ChartPath: "bitnami/mongodb"}
var mysqlDeployment = c.HelmConfig{ReleaseName: "mysql", ChartPath: "bitnami/mysql"}
var postgresqlDeployment = c.HelmConfig{ReleaseName: "postgresql", ChartPath: "bitnami/postgresql"}
var prometheusOperatorDeployment = c.HelmConfig{ReleaseName: "prometheus-operator", ChartPath: "prometheus-community/kube-prometheus-stack"}

// NewInstallCmd represents the install command
func NewDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the specified options to your Kubernetes cluster",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{
			"jaeger",
			"kafka",
			"mongodb",
			"mysql",
			"postgresql",
			"prometheus-operator",
		},
	}
	cmd.AddCommand(
		NewDeployJaegerCmd(),
		NewDeployKafkaCmd(),
		NewDeployMongodbCmd(),
		NewDeployMysqlCmd(),
		NewDeployPostgresqlCmd(),
		NewDeployPrometheusOperatorCmd(),
	)
	cmd.PersistentFlags().StringVarP(&chartVersion, "version", "v", "", "Override the Helm chart version for the specified deployments")
	cmd.PersistentFlags().StringVarP(&chartNamespace, "namespace", "n", "", "Override the namespace for the specified deployments")
	cmd.PersistentFlags().StringVar(&jaegerDeployment.Version, "jaeger-version", "", "Override the Jaeger Helm chart version")
	cmd.PersistentFlags().StringVar(&jaegerDeployment.Namespace, "jaeger-namespace", "observability", "Override the Jaeger namespace")
	cmd.PersistentFlags().StringVar(&kafkaDeployment.Version, "kafka-version", "", "Override the Kafka Helm chart version")
	cmd.PersistentFlags().StringVar(&kafkaDeployment.Namespace, "kafka-namespace", "kafka", "Override the Kafka namespace")
	cmd.PersistentFlags().StringVar(&mongoDeployment.Version, "mongodb-version", "", "Override the MongoDB Helm chart version")
	cmd.PersistentFlags().StringVar(&mongoDeployment.Namespace, "mongodb-namespace", "mongodb", "Override the MongoDB namespace")
	cmd.PersistentFlags().StringVar(&mysqlDeployment.Version, "mysql-version", "", "Override the MySQL Helm chart version")
	cmd.PersistentFlags().StringVar(&mysqlDeployment.Namespace, "mysql-namespace", "mysql", "Override the MySQL namespace")
	cmd.PersistentFlags().StringVar(&postgresqlDeployment.Version, "postgresql-version", "", "Override the PostgreSQL Helm chart version")
	cmd.PersistentFlags().StringVar(&postgresqlDeployment.Namespace, "postgresql-namespace", "postgresql", "Override the PostgreSQL namespace")
	cmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Version, "prom-op-version", "", "Override the Prometheus Operator Helm chart version")
	cmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Namespace, "prom-op-namespace", "monitoring", "Override the Prometheus Operator namespace")
	return cmd
}

func NewDeployJaegerCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "jaeger",
		Short: "Deploy Jaeger to your cluster",
		Run: func(cmd *cobra.Command, args []string) {
			func() {
				mapstructure.Decode(viper.GetStringMap("deploy.jaeger"), &jaegerDeployment)
				helm.Install(jaegerDeployment)
			}()
		},
	}
	return cmd
}

func NewDeployKafkaCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "kafka",
		Short: "Deploy Kafka to your cluster",
		Run: func(cmd *cobra.Command, args []string) {
			func() {
				mapstructure.Decode(viper.GetStringMap("deploy.kafka"), &kafkaDeployment)
				helm.Install(kafkaDeployment)
			}()
		},
	}
	return cmd
}

func NewDeployMongodbCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "mongodb",
		Short: "Deploy MongoDB to your cluster",
		Run: func(cmd *cobra.Command, args []string) {
			func() {
				mapstructure.Decode(viper.GetStringMap("deploy.mongodb"), &mongoDeployment)
				helm.Install(mongoDeployment)
			}()
		},
	}
	return cmd
}

func NewDeployMysqlCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "mysql",
		Short: "Deploy MySQL to your cluster",
		Run: func(cmd *cobra.Command, args []string) {
			func() {
				mapstructure.Decode(viper.GetStringMap("deploy.mysql"), &mysqlDeployment)
				helm.Install(mysqlDeployment)
			}()
		},
	}
	return cmd
}

func NewDeployPostgresqlCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "postgresql",
		Short: "Deploy PostgreSQL to your cluster",
		Run: func(cmd *cobra.Command, args []string) {
			func() {
				mapstructure.Decode(viper.GetStringMap("deploy.postgresql"), &postgresqlDeployment)
				helm.Install(postgresqlDeployment)
			}()
		},
	}
	return cmd
}

func NewDeployPrometheusOperatorCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "prometheus-operator",
		Short: "Deploy Prometheus Operator to your cluster",
		Run: func(cmd *cobra.Command, args []string) {
			func() {
				mapstructure.Decode(viper.GetStringMap("deploy.prometheus-operator"), &prometheusOperatorDeployment)
				helm.Install(prometheusOperatorDeployment)
			}()
		},
	}
	return cmd
}
