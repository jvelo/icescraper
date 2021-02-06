// Copyright 2021 Jérôme Velociter
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type IcecastServer struct {
	Url           string `yaml:"url"`
	SkipCertCheck bool   `yaml:"skipCertCheck"`
}

// Config is the top-level configuration for icecast monitor config file.
type Config struct {
	ScrapeInterval time.Duration   `yaml:"scrape_interval"`
	Servers        []IcecastServer `yaml:"icecast_servers"`
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
	}
	return cfg, nil
}

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}
	// Default to DefaultConfig for when YAML is empty
	*cfg = DefaultConfig

	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

var (
	// DefaultConfig is the default top-level configuration.
	DefaultConfig = Config{
		ScrapeInterval: 5 * time.Second,
	}
)
