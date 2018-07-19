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
  "fmt"
  "github.com/dotStart/Watchdog/cfg"
  "github.com/gin-gonic/gin"
  "html/template"
  "os"
  "path/filepath"
  "strings"
  "time"
)

func (s *Server) logRequest(ctx *gin.Context) {
  // Start timer
  start := time.Now()
  path := ctx.Request.URL.Path

  // Process request
  ctx.Next()

  // Log only when path is not being skipped
  end := time.Now()
  latency := end.Sub(start).Seconds() / 1000

  clientIP := ctx.ClientIP()
  method := ctx.Request.Method
  statusCode := ctx.Writer.Status()
  comment := ctx.Errors.ByType(gin.ErrorTypePrivate).String()

  s.logger.Debugf("[%.2f ms] %s - %s %d %s %s %s", latency, clientIP, ctx.Request.Proto, statusCode, method, path, comment)
}

func loadTemplates(files ...string) (*template.Template, error) {
  var t *template.Template = nil
  for _, f := range files {
    data, err := Asset(f)
    if err != nil {
      return nil, fmt.Errorf("failed to load template file %s: %s", f, err)
    }

    var tpl *template.Template
    if t == nil {
      t = template.New(filepath.Base(f))
      tpl = t
    } else {
      tpl = t.New(filepath.Base(f))
    }

    _, err = tpl.Parse(string(data))
    if err != nil {
      return nil, fmt.Errorf("failed to parse template file %s: %s", f, err)
    }
  }

  return t, nil
}

func (s *Server) getSite(ctx *gin.Context) *cfg.Site {
  host := ctx.Request.Host
  if i := strings.Index(host, ":"); i != -1 {
    host = host[:i]
  }

  c := s.cfgMgr.Get()
  site, ok := c.Sites[host]
  if !ok {
    s.logger.Debugf("requested site \"%s\" is not in site map", host)
    return nil
  }
  return &site
}

func (s *Server) getImagePath(site *cfg.Site, ext string) *string {
  dir := s.cfgMgr.Get().ImageDir
  if dir == nil {
    return nil
  }

  path := filepath.Join(*dir, site.DomainName+ext)
  _, err := os.Stat(path)
  if err != nil {
    return nil
  }
  return &path
}

func (s *Server) hasFavicon(site *cfg.Site) bool {
  return s.getImagePath(site, ".ico") != nil
}

func (s *Server) hasLogo(site *cfg.Site) bool {
  return s.getImagePath(site, ".png") != nil
}
