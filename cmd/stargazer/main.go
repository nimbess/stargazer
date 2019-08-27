// Copyright (c) 2019 Red Hat and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The main package for the Stargazer application.
package main

import (
	"context"
	"flag"
	"fmt"
	nimbessclientset "github.com/nimbess/stargazer/pkg/client/clientset/versioned"
	"github.com/nimbess/stargazer/pkg/config"
	"github.com/nimbess/stargazer/pkg/controller"
	unpv1 "github.com/nimbess/stargazer/pkg/crd/api/unp/v1"
	"github.com/nimbess/stargazer/pkg/etcdv3"
	log "github.com/sirupsen/logrus"
	extclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// VERSION is updated during the build process using git
var VERSION = "0.0.1"
var version = false
var cfgName = "stargazer"
var cfgPath = "."
var verbose = false

func init() {
	log.SetReportCaller(true)
	flag.BoolVar(&version, "v", version, "Display version")
	flag.StringVar(&cfgPath, "config-path", cfgPath, "Stargazer config file path")
	flag.StringVar(&cfgName, "config-name", cfgName, "Stargazer config file name")
	flag.BoolVar(&verbose, "verbose", verbose, "Enable debug logging")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	cfg := getConfig()

	// Get the context
	ctx := context.Background()

	// Get the k8s client
	k8sClient, err := getK8SClient(cfg.Kubeconfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to get k8s client api")
	}

	// Get the etcd client.
	etcdClient, err := getEtcdClient(cfg)
	if err != nil {
		log.WithError(err).Fatal("Failed to get etcd client")
	}

	// Register CRDs
	extClient, err := getK8sExtClient(cfg.Kubeconfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to get k8s extension client api")
	}
	if err := unpv1.CreateCRD(extClient); err != nil {
		log.WithError(err).Fatal("failed to create UNP CRD")
	}

	// setup the controller control structure
	stopCh := make(chan struct{})
	defer close(stopCh)

	controller.Run(cfg, k8sClient, etcdClient, ctx)

}

// getConfig gets the configuration
func getConfig() *config.Config {
	// Parse the user supplied config. If there are parsing errors then defaults will be returned.
	cfg := config.NewConfig()
	if err := cfg.Parse(cfgPath, cfgName); err != nil {
		log.WithField("cfgName", cfgName).WithField("cfgPath", cfgPath).
			WithError(err).Warn("Failed to parse config, using defaults")
	}

	log.WithField("config", fmt.Sprintf("%+v", *cfg)).Info("Configuration loaded")

	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)
	return cfg
}

// getK8SClient builds and returns a Kubernetes client.
func getK8SClient(kubeconfig string) (*nimbessclientset.Clientset, error) {
	// Build the kubeconfig.
	if kubeconfig == "" {
		log.Info("Using inClusterConfig")
	}
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %s", err)
	}

	// Get Kubernetes client.
	k8sClient, err := nimbessclientset.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubernetes client: %s", err)
	}

	return k8sClient, nil
}

// getK8sExtClient builds and returns a Kubernetes client.
func getK8sExtClient(kubeconfig string) (*extclientset.Clientset, error) {
	// Build the kubeconfig.
	if kubeconfig == "" {
		log.Info("Using inClusterConfig")
	}
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %s", err)
	}

	// Get Kubernetes client.
	k8sClient, err := extclientset.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubernetes client: %s", err)
	}

	return k8sClient, nil
}

func debugClient(k8sClient kubernetes.Interface) {
	log.Infof("k8sClient: %+v", k8sClient)

	pods, err := k8sClient.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Warnf("failed to get pods: %s", err)
	}
	log.Infof("pods: %+v", pods)

	pods, err = k8sClient.CoreV1().Pods("stargazer").List(metav1.ListOptions{})
	if err != nil {
		log.Warnf("failed to get default pods: %s", err)
	}
	log.Infof("pods: %+v", pods)

	nodes, err := k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Warnf("failed to get nodes: %s", err)
	}
	log.Infof("nodes: %+v", nodes)

	services, err := k8sClient.CoreV1().Services("stargazer").List(metav1.ListOptions{})
	if err != nil {
		log.Warnf("failed to get services: %s", err)
	}
	log.Infof("services: %+v", services)
}

func getEtcdClient(config *config.Config) (etcdv3.Client, error) {
	c, err := etcdv3.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get etcd client: %s", err)
	}
	return c, nil
}
