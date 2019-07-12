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
	"github.com/nimbess/stargazer/pkg/config"
	"github.com/nimbess/stargazer/pkg/controllers/controller"
	"github.com/nimbess/stargazer/pkg/controllers/node"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
	"time"
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

	// Get the k8s client api and shared informer factory.
	k8sClientset, err := getK8SClient(cfg.Kubeconfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to get k8s client api")
	}
	factory := getInformerFactory(k8sClientset, cfg)

	// TODO: ensure connection to etcd is available

	// setup the controller control structure
	ctx := context.Background()
	stopCh := make(chan struct{})
	defer close(stopCh)
	controllerCtrl := &controllerControl{
		ctx:            ctx,
		config:         cfg,
		stopCh:         stopCh,
		controllerInfo: make(map[string]*controllerInfo),
	}

	// Create an instance of each requested controller.
	// Store the instance along with it's number of workers into a manager list
	for _, controllerType := range strings.Split(cfg.Controllers, ",") {
		switch controllerType {
		case "node":
			nodeController := node.New(ctx, k8sClientset, cfg, factory)
			if nodeController == nil {
				log.WithField("controllerType", controllerType).Info("Failed to new controller")
				continue
			}
			controllerCtrl.controllerInfo["Node"] = &controllerInfo{
				controller: nodeController,
				workers:    cfg.NodeWorkers,
			}
		default:
			log.WithField("controllerType", controllerType).Info("Invalid controller")
		}
	}

	factory.Start(stopCh)
	controllerCtrl.RunControllers()
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
func getK8SClient(kubeconfig string) (kubernetes.Interface, error) {
	// Build the kubeconfig.
	if kubeconfig == "" {
		log.Info("Using inClusterConfig")
	}
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %s", err)
	}

	// Get Kubernetes clientset.
	k8sClientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubernetes clientset: %s", err)
	}

	return k8sClientset, nil
}

func debugClient(k8sClientset kubernetes.Interface) {
	log.Infof("k8sClientset: %+v", k8sClientset)

	pods, err := k8sClientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Warnf("failed to get pods: %s", err)
	}
	log.Infof("pods: %+v", pods)

	pods, err = k8sClientset.CoreV1().Pods("stargazer").List(metav1.ListOptions{})
	if err != nil {
		log.Warnf("failed to get default pods: %s", err)
	}
	log.Infof("pods: %+v", pods)

	nodes, err := k8sClientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Warnf("failed to get nodes: %s", err)
	}
	log.Infof("nodes: %+v", nodes)

	services, err := k8sClientset.CoreV1().Services("stargazer").List(metav1.ListOptions{})
	if err != nil {
		log.Warnf("failed to get services: %s", err)
	}
	log.Infof("services: %+v", services)
}

// getInformerFactory returns a SharedInformerFactory to use with the controllers.
func getInformerFactory(clientset kubernetes.Interface, cfg *config.Config) informers.SharedInformerFactory {
	// TODO: Use the config resync period
	return informers.NewSharedInformerFactory(clientset, time.Second*30)
}

// Object for keeping track of controller states and statuses.
type controllerControl struct {
	ctx            context.Context
	config         *config.Config
	stopCh         chan struct{}
	controllerInfo map[string]*controllerInfo
}

// Runs all the controllers and blocks indefinitely.
func (cc *controllerControl) RunControllers() {
	for controllerType, cs := range cc.controllerInfo {
		log.WithField("ControllerType", controllerType).Info("Starting controller")
		if err := cs.controller.Run(cs.workers, cc.stopCh); err != nil {
			log.WithField("ControllerType", controllerType).Warn("Failed to start controller")
		}
	}
	select {}
}

// Track controller information for each controller type.
type controllerInfo struct {
	controller controller.Controller
	workers    int
}
