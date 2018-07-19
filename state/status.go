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
)

type ServiceStatus uint8

const Unknown ServiceStatus = 0
const Operational ServiceStatus = 1
const Degraded ServiceStatus = 2
const Failed ServiceStatus = 3
const Maintenance ServiceStatus = 4

var statusNames = []string{"Unknown", "Operational", "Degraded", "Failed", "Maintenance"}

func (s *ServiceStatus) String() string {
  return statusNames[*s]
}

func convertConsulStatus(status string) ServiceStatus {
  switch status {
  case api.HealthPassing:
    return Operational
  case api.HealthWarning:
    return Degraded
  case api.HealthCritical:
    return Failed
  case api.HealthMaint:
    return Maintenance
  default:
    return Unknown
  }
}

func SelectStatus(a ServiceStatus, b ServiceStatus) ServiceStatus {
  if a > b {
    return a
  }
  return b
}

func OverallStatus(stateMap map[string]ServiceStatus) ServiceStatus {
  failCount := 0
  s := Unknown
  for _, state := range stateMap {
    if state == Unknown {
      return Unknown
    }
    if state == Failed {
      failCount++
    } else {
      s = SelectStatus(s, state)
    }
  }

  if failCount == len(stateMap) {
    return Failed
  }
  if failCount != 0 && s == Operational {
    return Degraded
  }
  return s
}
