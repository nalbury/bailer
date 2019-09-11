/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

//Create a kube client set
func kubeClient(kubeConfig string, clusterAuth bool) kubernetes.Clientset {
	var clientset *kubernetes.Clientset
	if clusterAuth {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			panic(err)
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}
	}
	return *clientset

}

//Faktory worker for running bailer jobs
func bail(ctx worker.Context, args ...interface{}) error {

	//This could be greatly improved if we can preserve the original structs passed to the worker...
	//For now we need to parse the values from the args (a slice of interface{}), and use type assertion

	//Get map of interfaces containing the alert labels
	alert := args[0].(map[string]interface{})
	alertLabels := alert["labels"].(map[string]interface{})

	//Get main bailer map of interfaces from the worker args
	bailer := args[1].(map[string]interface{})

	//Get map of interfaces containing the container image and tag
	bailerContainer := bailer["Container"].(map[string]interface{})

	//Get kubeconfig and cluster auth from the worker args
	kubeConfig := args[2].(string)
	clusterAuth := args[3].(bool)
	kubeClient := kubeClient(kubeConfig, clusterAuth)

	//stringified timestamp for unique naming
	ts := time.Now().Unix()
	var stamp string = fmt.Sprint(ts)

	//Get alert name from bailer map of interfaces and use it to create the jobName string
	alertName := strings.ToLower(bailer["Alert"].(string))
	jobName := "bailer-" + alertName + "-" + stamp
	jobNamespace := namespace

	//Get service account to run bailer job with from map of interfaces
	serviceAccountName := bailer["ServiceAccountName"].(string)

	//The image and tag cmd for the bailer job need to be cast to strings
	image := bailerContainer["Image"].(string) + ":" + bailerContainer["Tag"].(string)
	//cmd needs to be cast as a slice of strings
	var cmd []string
	for _, s := range bailer["Command"].([]interface{}) {
		cmd = append(cmd, s.(string))
	}

	//EnvVars from the labels on the alert, these can be used in bailer scripts
	var envVars []apiv1.EnvVar
	for key, value := range alertLabels {
		envVarKey := strings.ToUpper("ALERT_" + strings.Replace(key, "-", "_", -1))
		envVar := apiv1.EnvVar{Name: envVarKey, Value: value.(string)}
		envVars = append(envVars, envVar)

	}

	//Bailer job
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: jobNamespace,
		},

		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    alertName,
							Image:   image,
							Command: cmd,
							Env:     envVars,
						},
					},
					BackoffLimit:       0,
					RestartPolicy:      "Never",
					ServiceAccountName: serviceAccountName,
				},
			},
		},
	}
	//Run the bailer job
	fmt.Println("Creating job... ")

	jobsClient := kubeClient.BatchV1().Jobs(jobNamespace)
	result, err := jobsClient.Create(job)
	if err != nil {
		panic(fmt.Errorf("Unable to create job: %s \n", err))
	}
	fmt.Printf("Created job %q.\n", jobName)
	fmt.Printf("In Namespace %q.\n", jobNamespace)
	fmt.Printf("Job: %q \n", result)
	return nil
}

// faktoryCmd represents the faktory command
var faktoryCmd = &cobra.Command{
	Use:   "faktory",
	Short: "Background job processing for bailer",
	Long: `Bailer uses faktory to make upstream calls to kubernetes.
The Faktory server url can be set using the environment variable FAKTORY_URL, for example:
FAKTORY_URL=tcp://faktory.bailer.svc:7419
Defaults to localhost:7419`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr := worker.NewManager()
		mgr.Register("Bail", bail)
		mgr.Concurrency = 10
		mgr.ProcessStrictPriorityQueues("default")
		mgr.Run()
	},
}

func init() {
	rootCmd.AddCommand(faktoryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// faktoryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// faktoryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
