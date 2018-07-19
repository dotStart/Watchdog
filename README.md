Watchdog
========

A simple Consul powered status page.

Key Features
------------

* Multiple Sites on a single instance
* Multiple Services per Site
* Simple integration with (optimally remote) Consul instances

Prerequisites
-------------

* Golang 1.10 (or newer; required for building only)
* Consul 1.0 (or newer)

Installing from Source
----------------------

1. Execute `go get -d github.com/dotStart/Watchdog`
2. Execute `cd $GOPATH/src/github.com/dotStart/Watchdog`
3. Run `make`

Once the `make` command completes, you will find various executables (sorted by their
architecture and operating system) and their respective distribution packages inside of
the `build` directory.

Installing via Docker
---------------------

Watchdog is available as a Docker image from the
[GitLab registry](https://gitlab.com/dotStart/Watchdog/container_registry):

`docker pull registry.gitlab.com/dotstart/watchdog`

Note that the `latest` tag (which is also used when none is specified) includes unstable
releases of the application. It is highly recommended to select a specific stable version
if available instead.

Configuration
-------------

You may provide either a single configuration file or a directory of files via the
`-config` parameter. When multiple configuration files provide the same parameters,
previous instances of the same parameter will be overridden.

Example configuration files may be found in the [docs](docs) directory.

License
-------

```
Copyright [yyyy] [name of copyright owner]
and other copyright owners as documented in the project's IP log.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
