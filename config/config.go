// Copyright 2023 Quokka Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

type TLS struct {
	Cert string `toml:"cert"`
	Key  string `toml:"key"`
}

type SMTP struct {
	Address string `toml:"address"`
}

type IMAP struct {
	Address string `toml:"address"`
}

type IMAPS struct {
	Address string `toml:"address"`
}

type Submission struct {
	Address string `toml:"address"`
}

type Submissions struct {
	Address string `toml:"address"`
}

type Metrics struct {
	Address string `toml:"address"`
}

type Profiling struct {
	Address string `toml:"address"`
}

type Queue struct {
	Provider string `toml:"provider"`
}

type Auth struct {
}

type SPF struct {
}

type Log struct {
	Level string
}

type Config struct {
	Domain      string       `toml:"domain"`
	Log         *Log         `toml:"log"`
	Auth        *Auth        `toml:"auth"`
	IMAP        *IMAP        `toml:"imap"`
	IMAPS       *IMAPS       `toml:"imaps"`
	Metrics     *Metrics     `toml:"metrics"`
	Profiling   *Profiling   `toml:"profiling"`
	Queue       *Queue       `toml:"queue"`
	SMTP        *SMTP        `toml:"smtp"`
	SPF         *SPF         `toml:"spf"`
	Submission  *Submission  `toml:"submission"`
	Submissions *Submissions `toml:"submissions"`
	TLS         *TLS         `toml:"tls"`
}
