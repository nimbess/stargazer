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

package unp

import (
	"context"
	"github.com/nimbess/stargazer/pkg/config"
	unpv1 "github.com/nimbess/stargazer/pkg/crd/api/unp/v1"
	"github.com/nimbess/stargazer/pkg/etcdv3"
	"github.com/nimbess/stargazer/pkg/model"
	log "github.com/sirupsen/logrus"
	"path"
)

// Handler is implemented by any handler.
// The Handle method is used to process event
type UNP struct {
	etcdClient etcdv3.Client
	ctx        context.Context
}

// Init initializes handler configuration
func (u *UNP) Init(c *config.Config, etcdClient etcdv3.Client, ctx context.Context) error {
	u.etcdClient = etcdClient
	u.ctx = ctx
	return nil
}

// ObjectCreated creates entry in Nimbess DB with translated object
func (u *UNP) ObjectCreated(obj interface{}) {
	log.Infof("Created object found by controller: %v", obj)
	unpConf := obj.(*unpv1.UnifiedNetworkPolicy)
	kv, err := u.K8sToNimbess(unpConf)
	if err != nil {
		log.Errorf("Failed to convert K8S to Nimbess: %v", unpConf)
		return
	}
	err = u.etcdClient.Create(u.ctx, kv)
	if err != nil {
		log.Errorf("Failed to write to Nimbess etcd: %v", kv)
	}
}

// ObjectDeleted deletes entry in Nimbess DB with translated object
func (u *UNP) ObjectDeleted(name string) {
	log.Infof("Deleted object found by controller: %v", name)
	k := model.UNPKey{
		Name: name,
	}

	err := u.etcdClient.Delete(u.ctx, k)
	if err != nil {
		log.Errorf("Failed to delete key from Nimbess etcd: %v", k)
	}
}

// ObjectUpdated updates entry in Nimbess DB with translated object
func (u *UNP) ObjectUpdated(oldObj, newObj interface{}) {
	// TODO(trozet): implement
}

// TestHandler tests the handler configuration writing tests objects into DB
func (u *UNP) TestHandler() {

}

// K8stoNimbess translates a K8S UNP Config into a Nimbess Key/Value Pair to be written into ETCD
func (u *UNP) K8sToNimbess(unpConfig *unpv1.UnifiedNetworkPolicy) (*model.KVPair, error) {
	k := model.UNPKey{
		Name: path.Join(unpConfig.Namespace, unpConfig.Name),
	}

	kv := model.KVPair{Key: k, Value: unpConfig}

	log.WithFields(log.Fields{
		"k8s":    unpConfig,
		"KVPair": kv,
	}).Debug("Converted UNP")

	return &kv, nil

}
