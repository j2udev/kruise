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

var mongoDeployment = c.HelmConfig{ReleaseName: "mongodb", ChartPath: "bitnami/mongodb"}
var kafkaDeployment = c.HelmConfig{ReleaseName: "kafka", ChartPath: "bitnami/kafka"}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:       "deploy",
	Short:     "Deploy the specified options to your Kubernetes cluster",
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
					mapstructure.Decode(viper.GetStringMap("deploy.mongodb"), &mongoDeployment)
					helm.Install(mongoDeployment)
					wg.Done()
				}()
			case deployment == "kafka":
				go func() {
					mapstructure.Decode(viper.GetStringMap("deploy.kafka"), &kafkaDeployment)
					helm.Install(kafkaDeployment)
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
	deployCmd.PersistentFlags().StringVar(&kafkaDeployment.Version, "kafka-version", "", "Override the Kafka Helm chart version")
	deployCmd.PersistentFlags().StringVar(&kafkaDeployment.Namespace, "kafka-namespace", "kafka", "Override the Kafka namespace")
	deployCmd.PersistentFlags().StringVar(&mongoDeployment.Version, "mongodb-version", "", "Override the MongoDB Helm chart version")
	deployCmd.PersistentFlags().StringVar(&mongoDeployment.Namespace, "mongodb-namespace", "mongodb", "Override the MongoDB namespace")
}
