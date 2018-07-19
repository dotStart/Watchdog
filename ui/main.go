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
  "github.com/dotStart/Watchdog/cfg"
  "github.com/dotStart/Watchdog/state"
  "github.com/gin-gonic/gin"
  "github.com/op/go-logging"
  "path/filepath"
)

//go:generate go-bindata-assetfs -pkg ui -o data.gen.go -prefix resources/ resources/...

type Server struct {
  logger   *logging.Logger
  cfgMgr   *cfg.Manager
  stateMgr *state.Manager

  g *gin.Engine
}

func NewServer(addr string, uidir string, cfgMgr *cfg.Manager, stateMgr *state.Manager) (*Server, error) {
  s := &Server{
    logger:   logging.MustGetLogger("ui"),
    cfgMgr:   cfgMgr,
    stateMgr: stateMgr,

    g: gin.New(),
  }

  tpl, err := loadTemplates("template/index.html", "template/404.html")
  if err != nil {
    return nil, err
  }
  s.g.SetHTMLTemplate(tpl)
  s.g.Use(s.logRequest)

  if uidir == "" {
    s.g.StaticFS("/3rdParty/", assetFS())
    s.g.StaticFS("/scripts/", assetFS())
  } else {
    s.g.Static("/3rdParty/", filepath.Join(uidir, "3rdParty"))
    s.g.Static("/scripts/", filepath.Join(uidir, "scripts"))
  }

  s.g.GET("/", s.handleIndex)
  s.g.GET("/favicon.ico", s.handleFavicon)
  s.g.GET("/logo.png", s.handleLogo)
  s.g.NoRoute(s.handleNotFound)

  return s, nil
}

func (s *Server) Serve(addr ...string) error {
  return s.g.Run(addr...)
}
