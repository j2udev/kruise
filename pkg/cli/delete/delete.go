package cli

import (
	"sync"

	c "github.com/j2udevelopment/kruise/pkg/config"
	h "github.com/j2udevelopment/kruise/pkg/helm"
	u "github.com/j2udevelopment/kruise/pkg/utils"
	t "github.com/j2udevelopment/kruise/tpl"
	"github.com/spf13/cobra"
)

var chartNamespace string
var jaegerDeployment = c.HelmConfig{ReleaseName: "jaeger", ChartPath: "jaegertracing/jaeger"}
var kafkaDeployment = c.HelmConfig{ReleaseName: "kafka", ChartPath: "bitnami/kafka"}
var mongodbDeployment = c.HelmConfig{ReleaseName: "mongodb", ChartPath: "bitnami/mongodb"}
var mysqlDeployment = c.HelmConfig{ReleaseName: "mysql", ChartPath: "bitnami/mysql"}
var postgresqlDeployment = c.HelmConfig{ReleaseName: "postgresql", ChartPath: "bitnami/postgresql"}
var prometheusOperatorDeployment = c.HelmConfig{ReleaseName: "prometheus-operator", ChartPath: "prometheus-community/kube-prometheus-stack"}

// NewDeleteOpts returns options for the deploy command
func NewDeleteOpts() []u.Option {
	opts := []u.Option{
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
	return opts
}

// NewDeleteCmd represents the deploy command
func NewDeleteCmd(opts []u.Option) *cobra.Command {
	shallowDryRun := false
	cmd := &cobra.Command{
		Use:       "delete",
		Short:     "Delete the specified options from your Kubernetes cluster",
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
						h.Uninstall(shallowDryRun, &jaegerDeployment)
						wg.Done()
					}()
				case u.Contains(validArgsMap["kafka"], arg):
					go func() {
						h.Uninstall(shallowDryRun, &kafkaDeployment)
						wg.Done()
					}()
				case u.Contains(validArgsMap["mongodb"], arg):
					go func() {
						h.Uninstall(shallowDryRun, &mongodbDeployment)
						wg.Done()
					}()
				case u.Contains(validArgsMap["mysql"], arg):
					go func() {
						h.Uninstall(shallowDryRun, &mysqlDeployment)
						wg.Done()
					}()
				case u.Contains(validArgsMap["postgresql"], arg):
					go func() {
						h.Uninstall(shallowDryRun, &postgresqlDeployment)
						wg.Done()
					}()
				case u.Contains(validArgsMap["prometheus-operator"], arg):
					go func() {
						h.Uninstall(shallowDryRun, &prometheusOperatorDeployment)
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
	cmd.PersistentFlags().StringVarP(&chartNamespace, "namespace", "n", "", "Override the namespace for the specified deployments")
	cmd.PersistentFlags().StringVar(&jaegerDeployment.Namespace, "jaeger-namespace", "observability", "Override the Jaeger namespace")
	cmd.PersistentFlags().StringVar(&kafkaDeployment.Namespace, "kafka-namespace", "kafka", "Override the Kafka namespace")
	cmd.PersistentFlags().StringVar(&mongodbDeployment.Namespace, "mongodb-namespace", "mongodb", "Override the MongoDB namespace")
	cmd.PersistentFlags().StringVar(&mysqlDeployment.Namespace, "mysql-namespace", "mysql", "Override the MySQL namespace")
	cmd.PersistentFlags().StringVar(&postgresqlDeployment.Namespace, "postgresql-namespace", "postgresql", "Override the PostgreSQL namespace")
	cmd.PersistentFlags().StringVar(&prometheusOperatorDeployment.Namespace, "prom-op-namespace", "monitoring", "Override the Prometheus Operator namespace")

	return cmd
}
