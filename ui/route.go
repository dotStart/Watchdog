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
package ui

import (
  "errors"
  "github.com/dotStart/Watchdog/metadata"
  "github.com/dotStart/Watchdog/state"
  "github.com/gin-gonic/gin"
  "io/ioutil"
  "net/http"
)

func (s *Server) handleIndex(ctx *gin.Context) {
  site := s.getSite(ctx)
  if site == nil {
    s.handleNotFound(ctx)
    return
  }

  stateMap := s.stateMgr.GetServiceStatus(site)
  states := make([]state.ServiceStatus, len(stateMap))
  for _, stat := range stateMap {
    states = append(states, stat)
  }
  globalState := state.OverallStatus(states...)

  ctx.HTML(http.StatusOK, "index.html", &gin.H{
    "site": site,
    "images": &struct {
      Logo bool
    }{
      Logo: s.hasLogo(site),
    },
    "globalState": globalState,
    "state":       stateMap,
    "version":     metadata.VersionFull(),
  })
}

func (s *Server) handleFavicon(ctx *gin.Context) {
  site := s.getSite(ctx)
  if site == nil {
    s.handleNotFound(ctx)
    return
  }

  path := s.getImagePath(site, ".ico")
  if path == nil {
    s.handleNotFound(ctx)
    return
  }

  data, err := ioutil.ReadFile(*path)
  if err != nil {
    s.logger.Errorf("failed to load favicon for site \"%s\" at path %s: %s", site.DomainName, *path, err)
    ctx.Error(errors.New("cannot load favicon"))
    return
  }

  ctx.Data(http.StatusOK, "image/vnd.microsoft.icon", data)
}

func (s *Server) handleLogo(ctx *gin.Context) {
  site := s.getSite(ctx)
  if site == nil {
    s.handleNotFound(ctx)
    return
  }

  path := s.getImagePath(site, ".png")
  if path == nil {
    s.handleNotFound(ctx)
    return
  }

  data, err := ioutil.ReadFile(*path)
  if err != nil {
    s.logger.Errorf("failed to load logo for site \"%s\" at path %s: %s", site.DomainName, *path, err)
    ctx.Error(errors.New("cannot load favicon"))
    return
  }

  ctx.Data(http.StatusOK, "image/png", data)
}

func (s *Server) handleNotFound(ctx *gin.Context) {
  ctx.HTML(http.StatusNotFound, "404.html", nil)
}
