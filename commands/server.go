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
package commands

import (
  "context"
  "flag"
  "fmt"
  "github.com/dotStart/Watchdog/cfg"
  "github.com/dotStart/Watchdog/metadata"
  "github.com/dotStart/Watchdog/state"
  "github.com/dotStart/Watchdog/ui"
  "os"

  "github.com/google/subcommands"
  "github.com/op/go-logging"
  "golang.org/x/sync/errgroup"
)

type ServerCommand struct {
  baseCommand

  flagBindAddr string
  flagCfgPath  string
  flagUiDir    string
  flagLogLevel string
}

func (*ServerCommand) Name() string {
  return "serve"
}

func (*ServerCommand) Synopsis() string {
  return "starts serving requests on a specified address and port"
}

func (*ServerCommand) Usage() string {
  return `serve [options]

This command starts serving requests on a pre-configured address and port:

  $ watchdog serve

By default, the application will listen on all interfaces on port 46624. You may, however, also
specify a custom bind address:

  $ watchdog serve -bind=127.0.0.1:46625

Available Options:
`
}

func (c *ServerCommand) SetFlags(f *flag.FlagSet) {
  c.baseCommand.SetFlags(f)

  f.StringVar(&c.flagBindAddr, "bind", "", "specifies the address and port to bind to")
  f.StringVar(&c.flagCfgPath, "config", "", "specifies the location of the configuration file or directory")
  f.StringVar(&c.flagUiDir, "ui-dir", "", "specifies the location of the UI resources")
  f.StringVar(&c.flagLogLevel, "log-level", "info", "specifies the application log level")
}

func (c *ServerCommand) ParseEnvironmentVariables() {
  c.baseCommand.ParseEnvironmentVariables()

  if c.flagBindAddr == "" {
    bindAddr, set := os.LookupEnv("WATCHDOG_BIND_ADDR")
    if set {
      c.flagBindAddr = bindAddr
    } else {
      c.flagBindAddr = "0.0.0.0:46624"
    }
  }

  if c.flagCfgPath == "" {
    cfgPath, set := os.LookupEnv("WATCHDOG_CFG_PATH")
    if set {
      c.flagCfgPath = cfgPath
    } else {
      c.flagCfgPath = "config"
    }
  }
}

func (c *ServerCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
  c.ParseEnvironmentVariables()

  level, err := logging.LogLevel(c.flagLogLevel)
  if err != nil {
    fmt.Fprintf(os.Stderr, "illegal log level \"%s\": %s", c.flagLogLevel, err)
    return 1
  }

  if c.flagCfgPath == "" {
    fmt.Fprintf(os.Stderr, "-config parameter is required")
    return 1
  }

  fmt.Printf("==> Watchdog Configuration:\n\n")
  fmt.Printf("      Bind Address: %s\n", c.flagBindAddr)
  fmt.Printf("         Log Level: %s\n", c.flagLogLevel)
  fmt.Printf("           Version: %s\n", metadata.VersionFull())
  fmt.Printf("       Commit Hash: %s\n", metadata.CommitHash())
  fmt.Printf("   Build Timestamp: %s\n\n", metadata.BuildTimestamp())

  fmt.Printf("==> Starting Watchdog:\n\n")
  fmt.Printf("Application log will be streamed below:\n\n")

  format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} [%{level:.4s}] %{module} : %{color:reset} %{message}`)
  backend := logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))
  backend.SetLevel(level, "")
  logging.SetBackend(backend)

  cfgMgr, err := cfg.NewManager(c.flagCfgPath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to initialize server configuration: %s", err)
    return 2
  }

  stateMgr := state.New(cfgMgr)
  uiSrv, err := ui.NewServer(c.flagBindAddr, c.flagUiDir, cfgMgr, stateMgr)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to initialize server: %s", err)
    return 4
  }

  var g errgroup.Group
  g.Go(func() error {
    return uiSrv.Serve([]string{c.flagBindAddr}...)
  })
  g.Go(func() error {
    stateMgr.Tick()
    return nil
  })
  g.Go(func() error {
    cfgMgr.Hook()
    return nil
  })

  if err := g.Wait(); err != nil {
    fmt.Fprintf(os.Stderr, "failed to initialize server: %s", err)
    return 1
  }
  return 0
}
