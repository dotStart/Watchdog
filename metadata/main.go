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
package metadata

import "time"

var version = ""
var commitHash = ""
var buildTimestampRaw int64 = 0
var buildTimestamp = time.Unix(buildTimestampRaw, 0)

func Version() string {
  if version == "" {
    return "0.0.0"
  }

  return version
}

func VersionFull() string {
  if commitHash == "" {
    return Version() + "+dev"
  }

  return Version() + "+git-" + commitHash[:7]
}

func CommitHash() string {
  if commitHash == "" {
    return "dev"
  }

  return commitHash
}

func BuildTimestamp() time.Time {
  if buildTimestampRaw == 0 {
    return time.Now()
  }

  return buildTimestamp
}
