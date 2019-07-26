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

package handlers

import (
	"context"
	"github.com/nimbess/stargazer/pkg/config"
	"github.com/nimbess/stargazer/pkg/etcdv3"
	"github.com/nimbess/stargazer/pkg/handlers/unp"
)

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	Init(c *config.Config, etcdClient etcdv3.Client, ctx context.Context) error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
	TestHandler()
}

// Map maps each event handler function to a name for easily lookup
var Map = map[string]Handler{
	"default": &Default{},
	"UNP":     &unp.UNP{},
}

// Default handler implements Handler interface,
// print each event with JSON format
type Default struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (d *Default) Init(c *config.Config, etcdClient etcdv3.Client, ctx context.Context) error {
	return nil
}

// ObjectCreated sends events on object creation
func (d *Default) ObjectCreated(obj interface{}) {

}

// ObjectDeleted sends events on object deletion
func (d *Default) ObjectDeleted(obj interface{}) {

}

// ObjectUpdated sends events on object updation
func (d *Default) ObjectUpdated(oldObj, newObj interface{}) {

}

// TestHandler tests the handler configurarion by sending test messages.
func (d *Default) TestHandler() {

}
