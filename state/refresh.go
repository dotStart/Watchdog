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
  "github.com/hashicorp/consul/api"
  "strings"
)

func (mgr *Manager) Tick() {
  mgr.Refresh()

  for range mgr.refreshTicker.C {
    mgr.Refresh()
  }
}

func (mgr *Manager) Refresh() {
  mgr.logger.Infof("performing service refresh")
  c := mgr.cfgMgr.Get()
  statusMap := make(map[string]ServiceStatus)
  for _, site := range c.Sites {
    mgr.logger.Infof("checking state for site: %s", site.DomainName)
    for _, srv := range site.Services {
      client, err := api.NewClient(&api.Config{
        Address: srv.GetConsulAddress(),
      })

      if err != nil {
        mgr.logger.Errorf("  cannot connect to consul server at \"%s\": %s", srv.GetConsulAddress(), err)
        continue
      }

      state := Unknown
      if len(srv.ConsulServices) == 0 {
        mgr.logger.Errorf("  no configured consul services - skipped")
      } else {
        for _, srvName := range srv.ConsulServices {
          tag := ""

          i := strings.Index(srvName, "|")
          if i != -1 {
            tag = srvName[i+1:]
            srvName = srvName[:i]
          }

          mgr.logger.Debugf("  checking service \"%s\" with tag \"%s\"", srvName, tag)
          services, _, err := client.Health().Service(srvName, tag, false, &api.QueryOptions{})
          if err != nil {
            mgr.logger.Debugf("    server responded with an error: %s", err)
            state = SelectStatus(state, Unknown)
          } else {
            mgr.logger.Debugf("    found %d matches", len(services))
            var srvStatus ServiceStatus
            for _, s := range services {
              stat := convertConsulStatus(s.Checks.AggregatedStatus())
              mgr.logger.Debugf("    consul reports %s => %s", s.Checks.AggregatedStatus(), stat.String())

              srvStatus = SelectStatus(srvStatus, stat)
              mgr.logger.Debugf("    updated service state to %s", srvStatus.String())
            }

            state = SelectStatus(state, srvStatus)
            mgr.logger.Debugf("    => %s", srvStatus.String())
          }
        }
      }

      statusMap[getServiceKey(&site, &srv)] = state
    }
  }

  mgr.statusLock.Lock()
  defer mgr.statusLock.Unlock()
  mgr.statusMap = statusMap
}

func (mgr *Manager) getServiceList() []string {
  c := mgr.cfgMgr.Get()

  srvs := make([]string, 0)
  for _, site := range c.Sites {
    for _, srv := range site.Services {
      srvs = append(srvs, getServiceKey(&site, &srv))
    }
  }
  return srvs
}
