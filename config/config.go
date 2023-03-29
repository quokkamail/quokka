package config

type TLS struct {
	Cert string `toml:"cert"`
	Key  string `toml:"key"`
}

type Relay struct {
	Address string `toml:"address"`
}

type IMAP struct {
	Address string `toml:"address"`
}

type Submission struct {
	Address string `toml:"address"`
}

type Submissions struct {
	Address string `toml:"address"`
}

type Queue struct {
	Provider string `toml:"provider"`
}

type Auth struct {
	// RequireTLS bool `toml:"require_tls"`
}

type Config struct {
	Auth            *Auth        `toml:"auth"`
	IMAP            *IMAP        `toml:"imap"`
	Queue           *Queue       `toml:"queue"`
	SMTPRelay       *Relay       `toml:"smtp-relay"`
	SMTPSubmission  *Submission  `toml:"smtp-submission"`
	SMTPSubmissions *Submissions `toml:"smtp-submissions"`
	TLS             *TLS         `toml:"tls"`
}

var Default = Config{
	SMTPRelay: &Relay{
		Address: ":smtp",
	},
	IMAP: &IMAP{
		Address: ":imap",
	},
	SMTPSubmission: &Submission{
		Address: ":submission",
	},
	SMTPSubmissions: &Submissions{
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
