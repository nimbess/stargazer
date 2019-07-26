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
	typeUNP = reflect.TypeOf(UNP{})
)

type UNP struct {
	Namespace  string            `json:"namespace,omitempty"`
	Name       string            `json:"name,omitempty"`
	InternalIP string            `json:"internalip,omitempty"`
	Hostname   string            `json:"hostname,omitempty"`
	PodCIDR    string            `json:"podCIDR,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	UID        types.UID         `json:"uid,omitempty"`
}

type UNPKey struct {
	Name string
}

func (key UNPKey) defaultDeletePath() (string, error) {
	return key.defaultPath()
}

func (key UNPKey) defaultPath() (string, error) {
	if key.Name == "" {
		return "", errors.ErrorInsufficientIdentifiers{Name: "name"}
	}
	return fmt.Sprintf("/nimbess/unp/%s", key.Name), nil
}

func (key UNPKey) valueType() (reflect.Type, error) {
	return typeNode, nil
}

func (key UNPKey) String() string {
	return fmt.Sprintf("UNP(name=%s)", key.Name)
}

func (key UNPKey) KeyToDefaultDeletePath() (string, error) {
	return key.defaultPath()
}
