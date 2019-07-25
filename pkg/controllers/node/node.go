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
	"errors"
	"fmt"
	"github.com/nimbess/stargazer/pkg/config"
	"github.com/nimbess/stargazer/pkg/controllers/controller"
	"github.com/nimbess/stargazer/pkg/etcdv3"
	"github.com/nimbess/stargazer/pkg/model"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

// Controller implements the Controller interface. It is responsible for monitoring
// kubernetes nodes and responding to delete events by removing them from the Nimbess datastore.
type Controller struct {
	ctx        context.Context
	clientset  kubernetes.Interface
	config     *config.Config
	synced     cache.InformerSynced
	etcdClient etcdv3.Client
	lister     corelisters.NodeLister
	workqueue  workqueue.RateLimitingInterface
}

// New is the constructor for a Node Controller.
// config is assumed to have been validated earlier.
func New(
	ctx context.Context,
	clientset kubernetes.Interface,
	cfg *config.Config,
	etcdClient etcdv3.Client,
	factory informers.SharedInformerFactory) controller.Controller {
	if ctx == nil || clientset == nil || cfg == nil {
		return nil
	}

	informer := factory.Core().V1().Nodes().Informer()
	lister := factory.Core().V1().Nodes().Lister()
	c := &Controller{
		ctx:        ctx,
		clientset:  clientset,
		config:     cfg,
		synced:     informer.HasSynced,
		lister:     lister,
		etcdClient: etcdClient,
		workqueue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "node"),
	}

	//Setup event handlers.
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.handleObject(Add, obj)
		},
		UpdateFunc: func(old, new interface{}) {
			newNode := new.(*corev1.Node)
			oldNode := old.(*corev1.Node)
			if newNode == oldNode {
				return
			}
			c.handleObject(Update, new)
		},
		DeleteFunc: func(obj interface{}) {
			c.handleObject(Delete, obj)
		},
	})

	return c
}

// Run starts the node controller. It performs basic verifications and then launches
// worker threads.
func (c *Controller) Run(threadiness int, stopCh chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	log.Info("Starting controller")

	log.Info("Waiting for informer caches to sync")
	if !cache.WaitForCacheSync(stopCh, c.synced) {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	log.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	log.Info("Controller is now running")

	<-stopCh
	log.Info("Stopping Node controller")

	return nil
}

type Operation int

const (
	Add    Operation = 0
	Update Operation = 1
	Delete Operation = 2
)

type Event struct {
	key       string
	operation Operation
}

func (o Operation) String() string {
	names := [...]string{
		"Add",
		"Update",
		"Delete",
	}

	if o < Add || o > Delete {
		return "Unknown"
	}

	return names[o]
}

// handleObject will enqueue updates to the workqueue
func (c *Controller) handleObject(operation Operation, obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}

	c.workqueue.Add(Event{key, operation})
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var event Event
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if event, ok = obj.(Event); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		key = event.key
		operation := event.operation
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key, operation); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			// TODO: add maxRetries
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		log.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string, operation Operation) error {
	log.WithField("key", key).Debug("syncing ...")

	switch operation {
	case Add:
		return c.addNode(key)
	case Update:
		return c.addNode(key)
	case Delete:
		return c.deleteNode(key)
	default:
		return errors.New("Invalid operation")
	}
}

func (c *Controller) getK8sResource(key string) (*corev1.Node, error) {
	// Convert the namespace/name string into a distinct namespace and name
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil, nil
	}

	// Get the resource with this namespace/name
	k8sNode, err := c.lister.Get(name)
	if err != nil {
		// The resource may no longer exist, in which case we stop
		// processing.
		if apierrors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("resource '%s' in work queue no longer exists", key))
			return nil, nil
		}

		return nil, err
	}

	return k8sNode, nil
}

func (c *Controller) addNode(key string) error {
	k8sNode, err := c.getK8sResource(key)
	if err != nil {
		return err
	}
	if k8sNode == nil {
		log.Info("node no longer exists")
		return nil
	}

	nNode, _ := c.K8sNodeToNode(k8sNode)
	err = c.etcdClient.Create(c.ctx, nNode)
	return nil
}

func (c *Controller) deleteNode(key string) error {
	k8sNode, err := c.getK8sResource(key)
	if err != nil {
		return err
	}
	if k8sNode == nil {
		log.Info("node no longer exists")
		return nil
	}

	k := model.NodeKey{
		Hostname: k8sNode.Name,
	}
	err = c.etcdClient.Delete(c.ctx, &k)
	return nil
}

func (c *Controller) K8sNodeToNode(k8sNode *corev1.Node) (*model.KVPair, error) {
	v := model.Node{}
	v.Name = k8sNode.GetName()
	v.Namespace = k8sNode.GetNamespace()
	v.PodCIDR = k8sNode.Spec.PodCIDR
	v.UID = k8sNode.GetUID()
	v.Labels = k8sNode.GetLabels()

	// TODO: see if there is an interface that can be used like GetName instead of directly accessing structure
	for _, address := range k8sNode.Status.Addresses {
		switch address.Type {
		case corev1.NodeHostName:
			v.Hostname = address.Address
		case corev1.NodeInternalIP:
			v.InternalIP = address.Address
		}
	}
	k := model.NodeKey{
		Hostname: k8sNode.Name,
	}

	kv := model.KVPair{Key: k, Value: &v}

	log.WithFields(log.Fields{
		"k8s":    k8sNode,
		"KVPair": kv,
	}).Debug("Converted Node")

	return &kv, nil
}
