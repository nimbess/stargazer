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
	testName string
	cfgPath  string
	cfgName  string
}

type testoutput struct {
	error string
	cfg   config.Config
}

type testio struct {
	in       testinput
	expected testoutput
}

var defaultCfg = &config.Config{
	LogLevel: "Debug", Controllers: "node", NodeWorkers: 1, UNPWorkers: 1, Kubeconfig: "/etc/kubernetes/admin.conf"}

var tests = []testio{
	// pass: successful read and parse of config
	{testinput{"parse 1", "./testdata", "stargazer"}, testoutput{"success", *defaultCfg}},
	// pass: bad path, expected fail
	{testinput{"parse 2", "./badpath", "stargazer"}, testoutput{"badpath", *defaultCfg}},
	// pass: bad name, expected fail
	{testinput{"parse 3", "./testdata", "badname"}, testoutput{"badname", *defaultCfg}},
	// pass: successful read and parse of config with extra params that are ignored
	{testinput{"parse 4", "./testdata", "extra"}, testoutput{"success", *defaultCfg}},
	// pass: fail to read and parse of config with invalid params
	{testinput{"parse 5", "./testdata", "invalid"}, testoutput{"invalid", *defaultCfg}},
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

func fail(t *testing.T, testIo testio, in string, out string, got string) {
	_, file, line, _ := runtime.Caller(1)
	t.Errorf("%s\n%s:%d\nFor: %s\nExpected: %s\nGot: %s\n", testIo.in.testName, file, line, in, out, got)
	//t.Error("Test", testIo.in.testName, "line",fmt.Sprintf("%s:%d", file, line), "For", in, "Expected", out, "Got", got)
}

func TestConfig_Parse(t *testing.T) {
	for _, test := range tests {
		cfg := config.NewConfig()
		err := cfg.Parse(test.in.cfgPath, test.in.cfgName)
		if test.expected.error == "success" {
			if err == nil { // Parse passed
				if test.expected.cfg != *cfg { // fail if expected cfg not returned
					fail(t, test, fmt.Sprintf("%+v", test.in),
						fmt.Sprintf("%+v", test.expected.cfg),
						fmt.Sprintf("%+v", *cfg))
				}
			} else { // fail: error is unexpected
				fail(t, test, fmt.Sprintf("%+v", test.in),
					test.expected.error,
					fmt.Sprintf("%s", err))
			}
		} else {
			if err == nil { // fail if no Parse error
				fail(t, test, fmt.Sprintf("%+v", test.in),
					test.expected.error,
					"nil")
			}
		}
	}
}
