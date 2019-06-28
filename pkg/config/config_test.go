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

package config_test

import (
	"fmt"
	"github.com/nimbess/stargazer/pkg/config"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"testing"
)

type testinput struct {
	cfgPath string
	cfgName string
}

type testoutput struct {
	error string
	cfg   config.Config
}

type testio struct {
	in       testinput
	expected testoutput
}

var tests = []testio{
	// pass: successful read and parse of config
	{testinput{"./testdata", "stargazer"},
		testoutput{"success", config.Config{LogLevel: "Debug", NodeWorkers: 1, Controllers: "node"}}},
	// pass: bad path, expected fail
	{testinput{"./badpath", "stargazer"},
		testoutput{"bad path", config.Config{LogLevel: "Debug", NodeWorkers: 1, Controllers: "node"}}},
	// pass: bad name, expected fail
	{testinput{"./testdata", "badname"},
		testoutput{"bad name", config.Config{LogLevel: "Debug", NodeWorkers: 1, Controllers: "node"}}},
	// pass: successful read and parse of config with extra params that are ignored
	{testinput{"./testdata", "extra"},
		testoutput{"success", config.Config{LogLevel: "Debug", NodeWorkers: 1, Controllers: "node"}}},
	// pass: fail to read and parse of config with invalid params
	{testinput{"./testdata", "invalid"},
		testoutput{"invalid", config.Config{LogLevel: "Debug", NodeWorkers: 1, Controllers: "node"}}},
}

func init() {
	// Add file and line no to output logs
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			repopath := fmt.Sprintf("%s/src/github.com/nimbess/stargazer", os.Getenv("GOPATH"))
			filename := strings.Replace(f.File, repopath, "", -1)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		}})
}

func fail(t *testing.T, in string, out string, got string) {
	t.Error("For", in, "Expected", out, "Got", got)
}

func TestConfig_Parse(t *testing.T) {
	for _, test := range tests {
		cfg := config.NewConfig()
		err := cfg.Parse(test.in.cfgPath, test.in.cfgName)
		if test.expected.error == "success" {
			if err == nil { // Parse passed
				if test.expected.cfg != *cfg { // fail if expected cfg not returned
					fail(t, fmt.Sprintf("%+v", test.in),
						fmt.Sprintf("%+v", test.expected.cfg),
						fmt.Sprintf("%+v", *cfg))
				}
			} else { // fail: error is unexpected
				fail(t, fmt.Sprintf("%+v", test.in),
					test.expected.error,
					fmt.Sprintf("%s", err))
			}
		} else {
			if err == nil { // fail if no Parse error
				fail(t, fmt.Sprintf("%+v", test.in),
					test.expected.error,
					fmt.Sprintf("%s", err))
			}
		}
	}
}
