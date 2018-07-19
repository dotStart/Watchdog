/*
 * Copyright 2018 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cfg

import (
  "fmt"
  "github.com/hashicorp/hcl2/gohcl"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
  "time"

  "github.com/hashicorp/hcl2/hclparse"
)

const defaultConsulAddr = "http://localhost:8500"

type Config struct {
  RefreshRate    *time.Duration
  RawRefreshRate *string `hcl:"refresh-rate,attr"`
  ImageDir       *string `hcl:"image-dir,attr"`
  Sites          map[string]Site
  RawSites       []Site  `hcl:"site,block"`
}

func Empty() *Config {
  return &Config{
    Sites:    make(map[string]Site),
    RawSites: make([]Site, 0),
  }
}

func Load(path string) (*Config, error) {
  stat, err := os.Stat(path)
  if err != nil {
    return nil, err
  }

  if stat.IsDir() {
    return LoadDir(path)
  }
  return LoadFile(path)
}

func LoadFile(path string) (*Config, error) {
  parser := hclparse.NewParser()
  file, diag := parser.ParseHCLFile(path)
  if diag.HasErrors() {
    return nil, fmt.Errorf("failed to load configuration file \"%s\": %s", path, diag.Error())
  }

  cfg := Empty()
  diag = gohcl.DecodeBody(file.Body, nil, cfg)
  if diag.HasErrors() {
    return nil, fmt.Errorf("failed to parse configuration file \"%s\": %s", path, diag.Error())
  }
  err := cfg.parse()
  if err != nil {
    return nil, err
  }
  return cfg, nil
}

func LoadDir(path string) (*Config, error) {
  files, err := ioutil.ReadDir(path)
  if err != nil {
    return nil, fmt.Errorf("failed to load configuration directory \"%s\": %s", path, err)
  }

  cfg := Empty()
  for _, file := range files {
    if !strings.HasSuffix(file.Name(), ".hcl") {
      continue
    }

    c, err := LoadFile(filepath.Join(path, file.Name()))
    if err != nil {
      return nil, fmt.Errorf("failed to load configuration directory \"%s\": %s", path, err)
    }

    cfg.Merge(c)
  }

  return cfg, nil
}

func (c *Config) parse() error {
  if c.RawRefreshRate != nil {
    rate, err := time.ParseDuration(*c.RawRefreshRate)
    if err != nil {
      return err
    }

    c.RefreshRate = &rate
  }

  for _, site := range c.RawSites {
    c.Sites[site.DomainName] = site
  }

  return nil
}

func (c *Config) Merge(other *Config) {
  if other.RefreshRate != nil {
    c.RefreshRate = other.RefreshRate
  }

  if other.ImageDir != nil {
    c.ImageDir = other.ImageDir
  }

  for domainName, site := range other.Sites {
    c.Sites[domainName] = site
  }
}

type Site struct {
  DomainName  string    `hcl:"domain-name,label"`
  DisplayName string    `hcl:"display-name,attr"`
  Services    []Service `hcl:"service,block"`
}

type Service struct {
  Id             string   `hcl:"id,label"`
  DisplayName    string   `hcl:"display-name,attr"`
  ConsulAddr     *string  `hcl:"consul-addr,attr"`
  ConsulServices []string `hcl:"consul-services,attr"`
}

func (s *Service) GetConsulAddress() string {
  if s.ConsulAddr == nil {
    return defaultConsulAddr
  }

  return *s.ConsulAddr
}
