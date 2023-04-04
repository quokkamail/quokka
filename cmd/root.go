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

package cmd

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"

	"github.com/pelletier/go-toml/v2"
	"github.com/quokkamail/quokka/config"
	"github.com/quokkamail/quokka/smtp"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

type runRootOptions struct {
	config string
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "quokka",
		Short:         "Modern Mail Server",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger := slog.New(slog.NewTextHandler(os.Stderr))
			slog.SetDefault(logger)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			configFlag, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			if err := runRoot(runRootOptions{
				config: configFlag,
			}); err != nil {
				slog.Error(err.Error())
				os.Exit(1)
				return nil
			}

			return nil
		},
	}

	rootCmd.AddCommand(NewConfigCmd())

	rootCmd.Flags().StringP("config", "c", "config.toml", "configuration file")

	return rootCmd
}

func runRoot(opts runRootOptions) error {
	configBytes, err := os.ReadFile(opts.config)
	if err != nil {
		return err
	}

	var config config.Config
	if err := toml.Unmarshal(configBytes, &config); err != nil {
		return err
	}

	var tlsConfig *tls.Config
	if config.TLS != nil {
		slog.Info("initializing tls configuration")

		cert, err := tls.LoadX509KeyPair(config.TLS.Cert, config.TLS.Key)
		if err != nil {
			return fmt.Errorf("tls: %w", err)
		}

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
	}

	var smtpRelaySrv *smtp.Server
	if config.SMTPRelay != nil {
		slog.Info("starting smtp relay server", "address", config.SMTPRelay.Address)

		smtpRelaySrv = &smtp.Server{
			Address:                 config.SMTPRelay.Address,
			Domain:                  "quokka.local",
			TLSConfig:               tlsConfig,
			AuthenticationEncrypted: true,
		}

		go func() {
			if err := smtpRelaySrv.ListenAndServe(); err != nil {
				slog.Error(fmt.Errorf("smtp relay server: %w", err).Error())
			}
		}()
	}

	var smtpSubmissionSrv *smtp.Server
	if config.SMTPSubmission != nil {
		slog.Info("starting smtp submission server", "address", config.SMTPSubmission.Address)

		smtpSubmissionSrv = &smtp.Server{
			Address:                 config.SMTPSubmission.Address,
			Domain:                  "quokka.local",
			TLSConfig:               tlsConfig,
			AuthenticationEncrypted: true,
			AuthenticationMandatory: true,
		}

		go func() {
			if err := smtpSubmissionSrv.ListenAndServe(); err != nil {
				slog.Error(fmt.Errorf("smtp submission server: %w", err).Error())
			}
		}()
	}

	var smtpSubmissionsSrv *smtp.Server
	if config.SMTPSubmissions != nil {
		slog.Info("starting smtp submissions server", "address", config.SMTPSubmissions.Address)

		smtpSubmissionsSrv = &smtp.Server{
			Address:                 config.SMTPSubmissions.Address,
			Domain:                  "quokka.local",
			TLSConfig:               tlsConfig,
			AuthenticationEncrypted: true,
			AuthenticationMandatory: true,
		}

		go func() {
			if err := smtpSubmissionsSrv.ListenAndServe(); err != nil {
				slog.Error(fmt.Errorf("smtp submissions server: %w", err).Error())
			}
		}()
	}

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)

	<-interruptSignal
	slog.Info("got an interrupt signal")

	if smtpRelaySrv != nil {
		if err := smtpRelaySrv.Close(); err != nil {
			return err
		}
	}

	if smtpSubmissionSrv != nil {
		if err := smtpSubmissionSrv.Close(); err != nil {
			return err
		}
	}

	if smtpSubmissionsSrv != nil {
		if err := smtpSubmissionsSrv.Close(); err != nil {
			return err
		}
	}

	return nil
}
