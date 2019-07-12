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

package node_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/nimbess/stargazer/pkg/config"
	"github.com/nimbess/stargazer/pkg/controllers/controller"
	. "github.com/nimbess/stargazer/pkg/controllers/node"
	//"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	//"k8s.io/kubernetes/pkg/controller/testutil"
)

func TestNewController(t *testing.T) {
	cfg := config.NewConfig()
	clientset := k8sfake.NewSimpleClientset()
	factory := informers.NewSharedInformerFactory(clientset, 0)
	type args struct {
		ctx          context.Context
		k8sClientset kubernetes.Interface
		cfg          *config.Config
		informer     coreinformers.NodeInformer
	}
	tests := []struct {
		name string
		args args
		want controller.Controller
	}{
		{"test nil", args{nil, nil, nil, factory.Core().V1().Nodes()}, nil},
		{"test valid", args{nil, clientset, cfg, factory.Core().V1().Nodes()}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := New(test.args.ctx, test.args.k8sClientset, test.args.cfg, test.args.informer); !reflect.DeepEqual(got, test.want) {
				t.Errorf("New(%+v) = %+v\nwant: %+v", test.args, got, test.want)
			}
		})
	}
}
// TODO: fix the informer unit tests
/*
type fixture struct {
	t           *testing.T
	kubeclient  *k8sfake.Clientset
	//kubeobjects []runtime.Object
}

func newFixture(t *testing.T) *fixture {
	return &fixture{
		t:           t,
		kubeclient:  nil,
		//kubeobjects: []runtime.Object{},
	}
}

func (f *fixture) newController() (*controller.Controller, informers.SharedInformerFactory) {
	//f.kubeclient = k8sfake.NewSimpleClientset(&v1.NodeList{Items: []v1.Node{*testutil.NewNode("node")}})
	testNode := v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "machine1"}}
	fakenode := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "node1"},
		Status: v1.NodeStatus{
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceCPU):    resource.MustParse("10"),
				v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
			},
		},
	}
	//fakenode := *testutil.NewNode("node")
	//testNode := v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "machine1", UID: types.UID("machine1")}}
	clientset := k8sfake.NewSimpleClientset(&v1.NodeList{Items: []v1.Node{fakenode}})
	//k8sI := informers.NewSharedInformerFactory(f.kubeclient, 0)
	//c := New(nil, f.kubeclient, nil, k8sI.Core().V1().Nodes())
	//c.synced = true //TODO find way to access node.Controller fields compared to controller.Controller
	//return &c, k8sI
	return nil, nil
}

func (f *fixture) run(name string) {
	f.runController(name, true, false)
}

func (f *fixture) runController(name string, startInformers bool, expectedError bool) {
	c, k8sI := f.newController()
	if startInformers {
		stopCh := make(chan struct{})
		defer close(stopCh)
		k8sI.Start(stopCh)
	}
}
*/