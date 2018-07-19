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
package state

import (
  "github.com/dotStart/Watchdog/cfg"
  "github.com/op/go-logging"
  "sync"
  "time"
)

const defaultRefreshRate = time.Minute * 5

type Manager struct {
  logger *logging.Logger
  cfgMgr *cfg.Manager

  refreshTicker *time.Ticker
  statusLock    sync.RWMutex
  statusMap     map[string]ServiceStatus // FIXME: This is terrible
}

func New(cfgMgr *cfg.Manager) *Manager {
  c := cfgMgr.Get()
  refreshRate := defaultRefreshRate
  if c.RefreshRate != nil {
    refreshRate = *c.RefreshRate
  }

  return &Manager{
    logger: logging.MustGetLogger("state"),
    cfgMgr: cfgMgr,

    refreshTicker: time.NewTicker(refreshRate),
    statusMap:     make(map[string]ServiceStatus),
  }
}

func (mgr *Manager) GetServiceStatus(site *cfg.Site) map[string]ServiceStatus {
  mgr.statusLock.RLock()
  defer mgr.statusLock.RUnlock()

  status := make(map[string]ServiceStatus)
  for _, srv := range site.Services {
    s, ok := mgr.statusMap[getServiceKey(site, &srv)]
    if !ok {
      status[srv.Id] = Unknown
    } else {
      status[srv.Id] = s
    }
  }
  return status
}

func getServiceKey(site *cfg.Site, srv *cfg.Service) string {
  return site.DomainName + "#" + srv.Id
}
