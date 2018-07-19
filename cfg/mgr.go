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
  "github.com/op/go-logging"
  "os"
  "os/signal"
  "runtime"
  "sync"
  "syscall"
)

type Manager struct {
  logger *logging.Logger
  lock   sync.RWMutex
  path   string
  cfg    *Config
}

func NewManager(path string) (*Manager, error) {
  mgr := &Manager{
    logger: logging.MustGetLogger("cfg"),
    path:   path,
  }

  err := mgr.Reload()
  return mgr, err
}

func (mgr *Manager) Reload() error {
  mgr.lock.Lock()
  defer mgr.lock.Unlock()

  cfg, err := Load(mgr.path)
  if err != nil {
    return err
  }

  mgr.cfg = cfg
  return nil
}

func (mgr *Manager) Get() *Config {
  mgr.lock.RLock()
  defer mgr.lock.RUnlock()

  return mgr.cfg
}

func (mgr *Manager) Hook() {
  if runtime.GOOS == "windows" {
    mgr.logger.Warningf("signals are unavailable on windows - configuration file will not be reloaded")
    return
  }

  c := make(chan os.Signal, 1)
  signal.Notify(c, syscall.SIGHUP)

  for range c {
    mgr.logger.Infof("received SIGHUP - reloading application configuration")
    mgr.Reload()
  }
}
