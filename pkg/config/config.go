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

// Package config parses configuration files.
package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config stores the parsed configuration or defaults.
type Config struct {
	LogLevel    string
	Controllers string
	NodeWorkers int
	UNPWorkers  int
	Kubeconfig  string
	ResyncPeriod int64
}

// NewConfig is the constructor for Config.
func NewConfig() *Config {
	return &Config{
		LogLevel:    "info",
		Controllers: "node",
		NodeWorkers: 1,
		UNPWorkers:  1,
		Kubeconfig:  "",
		ResyncPeriod: 0,
	}
}

// Parse the configuration and store in Config.
// Defaults are returned if parsing fails.
func (c *Config) Parse(cfgPath string, cfgName string) error {
	vpr := viper.New()
	defaults := map[string]interface{}{
		"LogLevel":    c.LogLevel,
		"Controllers": c.Controllers,
		"NodeWorkers": c.NodeWorkers,
		"UNPWorkers":  c.UNPWorkers,
		"Kubeconfig":  c.Kubeconfig,
		"ResyncPeriod": c.ResyncPeriod,
	}
	for k, v := range defaults {
		vpr.SetDefault(k, v)
	}

	// TODO: tbd to use default paths and name
	vpr.AddConfigPath(cfgPath)
	vpr.SetConfigName(cfgName)
	vpr.AutomaticEnv()
	err := vpr.ReadInConfig()
	if err != nil {
		log.WithError(err).Warn("Failed to read config")
		return err
	}

	err = vpr.Unmarshal(c)
	if err != nil {
		log.WithError(err).Warn("Failed to unmarshal config")
	}
	return err
}
