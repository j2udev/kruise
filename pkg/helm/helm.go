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
package helm

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/j2udevelopment/kruise/pkg/config"
)

func constructChart(helmConfig *config.HelmConfig) {
	if helmConfig.ReleaseName == "" {
		log.Fatal("You must specify a Helm release name")
	}
	if helmConfig.ChartPath == "" {
		log.Fatal("You must specify a Helm chart")
	}
	if helmConfig.Namespace == "" {
		helmConfig.Namespace = "default"
	}
	if len(helmConfig.Args) == 0 {
		helmConfig.Args = []string{
			"upgrade", "-i",
			helmConfig.ReleaseName,
			helmConfig.ChartPath,
			"--namespace", helmConfig.Namespace,
			"--create-namespace",
		}
	}
	if helmConfig.Version != "" {
		helmConfig.Args = append(helmConfig.Args, "--version", helmConfig.Version)
	}
}

func Install(helmConfig config.HelmConfig) {
	helmCheck := exec.Command("command", "-v", "helm")
	if err := helmCheck.Run(); err != nil {
		log.Fatal("Helm does not appear to be installed")
	}
	constructChart(&helmConfig)
	for _, val := range helmConfig.Values {
		helmConfig.Args = append(helmConfig.Args, "-f", val)
	}
	for _, val := range helmConfig.ExtraArgs {
		helmConfig.Args = append(helmConfig.Args, val)
	}
	cmd := exec.Command("helm", helmConfig.Args...)
	fmt.Println(cmd)
	fmt.Println("Attempting to deploy " + helmConfig.ChartPath)
	if err := cmd.Run(); err != nil {
		log.Fatal("Something went wrong, is Kubernetes running?")
	} else {
		fmt.Println("Successfully deployed " + helmConfig.ChartPath)
	}
}

func Uninstall(helmConfig config.HelmConfig) {
	helmCheck := exec.Command("command", "-v", "helm")
	if err := helmCheck.Run(); err != nil {
		log.Fatal("Helm does not appear to be installed")
	}
	constructChart(&helmConfig)
	deleteArgs := []string{
		"uninstall",
		helmConfig.ReleaseName,
		"--namespace",
		helmConfig.Namespace,
	}
	for _, val := range helmConfig.ExtraArgs {
		deleteArgs = append(helmConfig.Args, val)
	}
	cmd := exec.Command("helm", deleteArgs...)
	fmt.Println(cmd)
	fmt.Println("Attempting to delete " + helmConfig.ChartPath)
	if err := cmd.Run(); err != nil {
		log.Fatal("Something went wrong, is Kubernetes running?")
	} else {
		fmt.Println("Successfully deleted " + helmConfig.ChartPath)
	}
}
