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
	"time"
)

// Controllers holds the enabled/disabled controller types
type Controllers struct {
	UNP bool
}

// Config stores the parsed configuration or defaults.
type Config struct {
	LogLevel        string
	Controllers     Controllers
	Kubeconfig      string
	ResyncPeriod    int64
	EtcdEndpoints   string
	EtcdDialTimeout time.Duration
}

// NewConfig is the constructor for Config.
func NewConfig() *Config {
	ctrl := Controllers{UNP: true}
	return &Config{
		LogLevel:        "info",
		Controllers:     ctrl,
		Kubeconfig:      "",
		ResyncPeriod:    0,
		EtcdEndpoints:   "http://127.0.0.1:52379",
		EtcdDialTimeout: 1 * time.Second,
	}
}

// Parse the configuration and store in Config.
// Defaults are returned if parsing fails.
func (c *Config) Parse(cfgPath string, cfgName string) error {
	vpr := viper.New()
	defaults := map[string]interface{}{
		"LogLevel":        c.LogLevel,
		"Controllers":     c.Controllers,
		"Kubeconfig":      c.Kubeconfig,
		"ResyncPeriod":    c.ResyncPeriod,
		"EtcdEndpoints":   c.EtcdEndpoints,
		"EtcdDialTimeout": c.EtcdDialTimeout,
	}
	for k, v := range defaults {
		vpr.SetDefault(k, v)
	}

	// TODO: tbd to use default paths and name
	vpr.AddConfigPath(cfgPath)
	vpr.SetConfigName(cfgName)
	vpr.AutomaticEnv()
	_ = vpr.BindEnv("EtcdEndpoints", "ETCDCTL_ENDPOINTS")
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
