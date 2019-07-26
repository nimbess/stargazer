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

package controller

import (
	"context"
	"fmt"
	nimbessclientset "github.com/nimbess/stargazer/pkg/client/clientset/versioned"
	unpinformer "github.com/nimbess/stargazer/pkg/client/informers/externalversions"
	"github.com/nimbess/stargazer/pkg/config"
	"github.com/nimbess/stargazer/pkg/etcdv3"
	"github.com/nimbess/stargazer/pkg/handlers"
	"github.com/nimbess/stargazer/pkg/signals"
	"github.com/nimbess/stargazer/pkg/utils"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const maxRetries = 5

var serverStartTime time.Time

// Event indicate the informerEvent
type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

// Controller object
type Controller struct {
	logger       *log.Entry
	clientset    nimbessclientset.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler handlers.Handler
	etcdClient   etcdv3.Client
}

// Runs stargazer and then waits for process termination signals
func Run(conf *config.Config, kubeClient *nimbessclientset.Clientset, etcdClient etcdv3.Client, ctx context.Context) {
	v := reflect.ValueOf(conf.Controllers)
	defer utilruntime.HandleCrash()
	stopCh := signals.SetupSignalHandler()
	ctrlType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == true {
			log.Infof("Enabling controller for: %s", ctrlType.Field(i).Name)
			// create handler
			thisHandler, ok := handlers.Map[ctrlType.Field(i).Name]
			if !ok {
				log.Fatalf("Unsupported handler for controller: %s", ctrlType.Field(i).Name)
			}
			if err := thisHandler.Init(conf, etcdClient, ctx); err != nil {
				log.Fatalf("Failed to init handler: %s", ctrlType.Field(i).Name)
			}
			c := Start(conf, kubeClient, thisHandler, stopCh)
			defer c.queue.ShutDown()
			log.Infof("Controller started: %s", ctrlType.Field(i).Name)
		}
	}

	<-stopCh
}

// Start prepares a watcher and run corresponding controllers. Non-blocking. Returns new controller object.
func Start(conf *config.Config, kubeClient *nimbessclientset.Clientset, eventHandler handlers.Handler,
	stopCh <-chan struct{}) *Controller {

	var informer cache.SharedIndexInformer
	var resType string

	nimbessInformerFactory := unpinformer.NewSharedInformerFactory(kubeClient, time.Second*30)

	if conf.Controllers.UNP {
		informer = nimbessInformerFactory.Nimbess().V1().UnpConfigs().Informer()
		resType = "unp"
	}
	c := newResourceController(kubeClient, eventHandler, informer, resType)

	if err := c.Run(stopCh); err != nil {
		log.Fatalf("Error running controller: %s, error: %v", resType, err)
	}
	return c
}

func newResourceController(client *nimbessclientset.Clientset, eventHandler handlers.Handler, informer cache.SharedIndexInformer, resourceType string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var newEvent Event
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				utilruntime.HandleError(err)
			}
			newEvent.eventType = "create"
			newEvent.resourceType = resourceType
			log.WithField("pkg", "stargazer-"+resourceType).Infof("Processing add to %v: %s", resourceType, newEvent.key)
			queue.Add(newEvent)
		},
		UpdateFunc: func(old, new interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(old)
			newEvent.eventType = "update"
			newEvent.resourceType = resourceType
			log.WithField("pkg", "stargazer-"+resourceType).Infof("Processing update to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
		DeleteFunc: func(obj interface{}) {
			newEvent.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			newEvent.eventType = "delete"
			newEvent.resourceType = resourceType
			newEvent.namespace = utils.GetObjectMetaData(obj).Namespace
			log.WithField("pkg", "stargazer-"+resourceType).Infof("Processing delete to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
	})

	return &Controller{
		logger:       log.WithField("pkg", "stargazer-"+resourceType),
		clientset:    client,
		informer:     informer,
		queue:        queue,
		eventHandler: eventHandler,
	}
}

// Run starts the stargazer controller
func (c *Controller) Run(stopCh <-chan struct{}) error {

	c.logger.Info("Starting stargazer controller")
	//serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return fmt.Errorf("failed to wait for caches to sync")
	}

	c.logger.Info("Stargazer controller synced and ready")

	go wait.Until(c.runWorker, time.Second, stopCh)
	c.logger.Info("Stargazer controller started")
	return nil
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()
	if quit {
		c.logger.Info("queue shutdown")
		return false
	}
	c.logger.Debugf("processing new item %v", newEvent)
	defer c.queue.Done(newEvent)
	err := c.processItem(newEvent.(Event))
	c.logger.Debugf("Done processing item, err is %v", err)
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(newEvent)
	} else if c.queue.NumRequeues(newEvent) < maxRetries {
		c.logger.Errorf("Error processing %s (will retry): %v", newEvent.(Event).key, err)
		c.queue.AddRateLimited(newEvent)
	} else {
		// err != nil and too many retries
		c.logger.Errorf("Error processing %s (giving up): %v", newEvent.(Event).key, err)
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(newEvent Event) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(newEvent.key)
	if err != nil {
		return fmt.Errorf("error fetching object with key %s from store: %v", newEvent.key, err)
	}
	// get object's metadata
	//objectMeta := utils.GetObjectMetaData(obj)

	// process events based on its type
	switch newEvent.eventType {
	case "create":
		// compare CreationTimestamp and serverStartTime and alert only on latest events
		// Could be Replaced by using Delta or DeltaFIFO
		//if objectMeta.CreationTimestamp.Sub(serverStartTime).Seconds() > 0 {
		c.logger.Debug("Calling create handler")
		c.eventHandler.ObjectCreated(obj)
		return nil
		//}
	case "update":
		c.logger.Debug("Inside update handler")
		/**
		kbEvent := event.Event{
			Kind: newEvent.resourceType,
			Name: newEvent.key,
		}
		c.eventHandler.ObjectUpdated(obj, kbEvent)
		*/
		return nil
	case "delete":
		c.logger.Debug("Inside delete handler")
		c.eventHandler.ObjectDeleted(obj)
		return nil
	}
	return nil
}
