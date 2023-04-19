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

type Pprof struct {
	Address string `toml:"address"`
}

type Queue struct {
	Provider string `toml:"provider"`
}

type Auth struct {
	// RequireTLS bool `toml:"require_tls"`
}

type Config struct {
	Auth        *Auth        `toml:"auth"`
	IMAP        *IMAP        `toml:"imap"`
	IMAPS       *IMAPS       `toml:"imaps"`
	Metrics     *Metrics     `toml:"metrics"`
	Pprof       *Pprof       `toml:"pprof"`
	Queue       *Queue       `toml:"queue"`
	SMTP        *SMTP        `toml:"smtp"`
	Submission  *Submission  `toml:"submission"`
	Submissions *Submissions `toml:"submissions"`
	TLS         *TLS         `toml:"tls"`
}

var Default = Config{
	SMTP: &SMTP{
		Address: ":smtp",
	},
	IMAP: &IMAP{
		Address: ":imap",
	},
	Submission: &Submission{
		Address: ":submission",
	},
	Submissions: &Submissions{
		Address: ":465",
	},
	TLS: &TLS{
		Cert: "cert.pem",
		Key:  "key.pem",
	},
	Auth: &Auth{
		// RequireTLS: true,
	},
	Queue: &Queue{
		Provider: "inmemory",
	},
}
