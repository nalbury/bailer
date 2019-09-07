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
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

// Vars for cmd line flag values
var cfgFile string
var kubeConfig string
var namespace string
var clusterAuth bool

// PAYLOAD FROM ALERTMANAGER

//Main type for incoming alert payloads
type PrometheusWebhook struct {
	Status string   `json:"status"`
	Alerts []Alerts `json:"alerts"`
}

//Labels for each alert
type Labels struct {
	Alertname string `json:"alertname"`
	Namespace string `json:"namespace"`
	PodName   string `json:"pod"`
	Severity  string `json:"severity"`
}

//Alert type
type Alerts struct {
	Status   string    `json:"status"`
	Labels   Labels    `json:"labels"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}

//BAILER CONFIG FILE

//Container image and tag that will be used by bailer
type Container struct {
	Image string `mapstructure:"image"`
	Tag   string `mapstructure:"tag"`
}

//Bailer type

type Bailer struct {
	Alert              string            `mapstructure:"alert"`
	Labels             map[string]string `mapstructure:"labels"`
	Command            []string          `mapstructure:"command"`
	Container          Container         `mapstructure:"container"`
	ServiceAccountName string            `mapstructure:"serviceAccountName"`
}

//Slice of Bailer structs
type Bailers struct {
	Bailers []Bailer `mapstructure:"bailers"`
}

//Check if the namespace and pod label filters match the labels on the alert
//Pod filter compiles to a reg expression and we attempt to find a match anywhere
//in the PodName label from alertmanager
func needsBailing(labels Labels, bailer Bailer) bool {

	pod, _ := regexp.Compile(bailer.Labels["pod"])
	return (labels.Alertname == bailer.Alert &&
		labels.Namespace == bailer.Labels["namespace"] &&
		pod.MatchString(labels.PodName))
}

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

//Run the bailer
func bail(labels Labels, bailer Bailer, kubeClient kubernetes.Clientset) {
	//Create string from unix ts for unique job name
	ts := time.Now().Unix()
	var stamp string = fmt.Sprint(ts)

	alert_name := strings.ToLower(bailer.Alert)
	jobName := "bailer-" + alert_name + "-" + stamp
	jobNamespace := namespace
	serviceAccountName := bailer.ServiceAccountName

	//The image and cmd for the bailer job
	image := bailer.Container.Image + ":" + bailer.Container.Tag
	cmd := bailer.Command

	//EnvVars from the labels on the alert, these can be used in bailer scripts
	envVars := []apiv1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: labels.PodName,
		},
		{
			Name:  "NAMESPACE",
			Value: labels.Namespace,
		},
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
							Name:    alert_name,
							Image:   image,
							Command: cmd,
							Env:     envVars,
						},
					},
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
	fmt.Printf("Started At: %q \n", result.Status.StartTime)
}

var rootCmd = &cobra.Command{
	Use:   "bailer",
	Short: "Bail out your kubernetes cluster",
	Long: `Bail out a leaky kubernetes cluster.
Bailer provides a webhook driven way of triggering kubernetes jobs from prometheus alerting rules.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Unmarshal config into []Bailers
		var bailers Bailers

		err := viper.Unmarshal(&bailers)
		if err != nil {
			panic(fmt.Errorf("unable to decode into struct, %v \n", err))
		}
		fmt.Println(bailers)
		fmt.Println(kubeConfig)

		// Create kube client
		kc := kubeClient(kubeConfig, clusterAuth)

		// Start jobrunner for background jobs
		jobrunner.Start()
		// Start up gin
		router := gin.Default()

		alert := router.Group("/alert")
		{
			alert.POST("/", func(c *gin.Context) {
				//Unmarshal payload into a PrometheusWebhook
				var p PrometheusWebhook
				c.BindJSON(&p)
				//For each alert in the the webhook payload check if the alert is firing
				for _, a := range p.Alerts {
					if a.Status == "firing" {
						//If it's a firing alert, iterate through the Bailers in the config
						//and use needsBailing to check if the alert matches a configured Bailer
						for _, b := range bailers.Bailers {
							if needsBailing(a.Labels, b) {
								//If the alert matches a bailer config:
								// - Use kube client to create a kube job
								// - Job should use bailer image spec
								// - Job should user bailer cmd spec
								// - Job should have alert labels available as env vars
								fmt.Println("NAMESPACE:", a.Labels.Namespace)
								fmt.Println("POD_NAME:", a.Labels.PodName)
								fmt.Println("Will use container image:", b.Container.Image)
								fmt.Println("To run:", strings.Join(b.Command, " "))
								bail(a.Labels, b, kc)
							}
						}
					}
				}
				c.JSON(http.StatusOK, gin.H{})
			})
		}
		router.Run(":3000")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		filepath.Join(usr.HomeDir, ".bailer", "config.yaml"),
		"bailer config file (default is $HOME/.bailer/config.yaml)",
	)

	rootCmd.PersistentFlags().StringVar(
		&kubeConfig,
		"kubeconfig",
		filepath.Join(usr.HomeDir, ".kube", "config"),
		"kube config file (default is $HOME/.kube/config)",
	)

	rootCmd.PersistentFlags().StringVar(
		&namespace,
		"namespace",
		"",
		"namespace that bailer will run jobs in",
	)

	rootCmd.PersistentFlags().BoolVar(
		&clusterAuth,
		"clusterauth",
		false,
		"use RBAC auth",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bailer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bailer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
