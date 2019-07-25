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

package model

import (
	"fmt"
	"github.com/nimbess/stargazer/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
)

var (
	typeNode = reflect.TypeOf(Node{})
)

type Node struct {
	Namespace  string            `json:"namespace,omitempty"`
	Name       string            `json:"name,omitempty"`
	InternalIP string            `json:"internalip,omitempty"`
	Hostname   string            `json:"hostname,omitempty"`
	PodCIDR    string            `json:"podCIDR,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	UID        types.UID         `json:"uid,omitempty"`
}

type NodeKey struct {
	Hostname string
}

func (key NodeKey) defaultDeletePath() (string, error) {
	return key.defaultPath()
}

func (key NodeKey) defaultPath() (string, error) {
	if key.Hostname == "" {
		return "", errors.ErrorInsufficientIdentifiers{Name: "name"}
	}
	return fmt.Sprintf("/nimbess/host/%s", key.Hostname), nil
}

func (key NodeKey) valueType() (reflect.Type, error) {
	return typeNode, nil
}

func (key NodeKey) String() string {
	return fmt.Sprintf("Node(name=%s)", key.Hostname)
}

func (key NodeKey) KeyToDefaultDeletePath() (string, error) {
	return key.defaultPath()
}
