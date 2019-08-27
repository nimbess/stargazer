/*
Copyright The Kubernetes Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	versioned "github.com/nimbess/stargazer/pkg/client/clientset/versioned"
	internalinterfaces "github.com/nimbess/stargazer/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/nimbess/stargazer/pkg/client/listers/unp/v1"
	unpv1 "github.com/nimbess/stargazer/pkg/crd/api/unp/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// UnifiedNetworkPolicyInformer provides access to a shared informer and lister for
// UnifiedNetworkPolicies.
type UnifiedNetworkPolicyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.UnifiedNetworkPolicyLister
}

type unifiedNetworkPolicyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewUnifiedNetworkPolicyInformer constructs a new informer for UnifiedNetworkPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewUnifiedNetworkPolicyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredUnifiedNetworkPolicyInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredUnifiedNetworkPolicyInformer constructs a new informer for UnifiedNetworkPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredUnifiedNetworkPolicyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NimbessV1().UnifiedNetworkPolicies(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NimbessV1().UnifiedNetworkPolicies(namespace).Watch(options)
			},
		},
		&unpv1.UnifiedNetworkPolicy{},
		resyncPeriod,
		indexers,
	)
}

func (f *unifiedNetworkPolicyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredUnifiedNetworkPolicyInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *unifiedNetworkPolicyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&unpv1.UnifiedNetworkPolicy{}, f.defaultInformer)
}

func (f *unifiedNetworkPolicyInformer) Lister() v1.UnifiedNetworkPolicyLister {
	return v1.NewUnifiedNetworkPolicyLister(f.Informer().GetIndexer())
}
