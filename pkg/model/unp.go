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
	"github.com/nimbess/stargazer/pkg/crd/api/unp/v1"
	"reflect"
)

var (
	typeUNP = reflect.TypeOf(v1.UnifiedNetworkPolicy{})
)

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
	return typeUNP, nil
}

func (key UNPKey) String() string {
	return fmt.Sprintf("UNP(name=%s)", key.Name)
}

func (key UNPKey) KeyToDefaultDeletePath() (string, error) {
	return key.defaultPath()
}
