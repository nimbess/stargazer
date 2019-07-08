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

// Package node is the Node controller implementation.
package node

import (
	"context"
	"github.com/nimbess/stargazer/pkg/config"
	"github.com/nimbess/stargazer/pkg/controller"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"sync"
)

// Controller implements the Controller interface. It is responsible for monitoring
// kubernetes nodes and responding to delete events by removing them from the Nimbess datastore.
type Controller struct {
	ctx         context.Context
	informer    cache.Controller
	indexer     cache.Indexer
	config      *config.Config
	nodemapper  map[string]string
	nodemapLock sync.Mutex
}

// NewController is the constructor for a Node Controller.
// config is assumed to have been validated earlier.
func NewController(ctx context.Context, cfg *config.Config) controller.Controller {
	return &Controller{
		ctx:        ctx,
		nodemapper: map[string]string{},
		config:     cfg,
	}
}

// Run starts the node controller. It performs basic verifications and then launches
// worker threads.
func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	log.Info("Starting Node controller")

	log.Info("Waiting to sync with Kubernetes API (Nodes)")

	log.Info("Finished syncing with Kubernetes API (Nodes)")

	log.Info("Node controller is now running")

	<-stopCh
	log.Info("Stopping Node controller")
}
