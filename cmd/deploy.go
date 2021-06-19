/*
Copyright © 2021 Joshua Ward <j2udevelopment@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
