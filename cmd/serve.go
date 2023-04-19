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
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"

	"github.com/pelletier/go-toml/v2"
	"github.com/quokkamail/quokka/config"
	"github.com/quokkamail/quokka/metrics"
	"github.com/quokkamail/quokka/pprof"
	"github.com/quokkamail/quokka/smtp"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

type runServeOptions struct {
	config string
}

func NewServeCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "serve",
		Short: "Parse the configuration file and start the server(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFlag, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			if err := runServe(runServeOptions{
				config: configFlag,
			}); err != nil {
				slog.Error(err.Error())
				os.Exit(1)
				return nil
			}

			return nil
		},
	}

	rootCmd.Flags().StringP("config", "c", "config.toml", "configuration file")

	return rootCmd
}

func runServe(opts runServeOptions) error {
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

	var smtpSrv *smtp.Server
	if config.SMTP != nil {
		slog.Info("starting smtp server", "address", config.SMTP.Address)

		smtpSrv = &smtp.Server{
			Address:                 config.SMTP.Address,
			Domain:                  "quokka.local",
			TLSConfig:               tlsConfig,
			AuthenticationEncrypted: true,
		}

		go func() {
			if err := smtpSrv.ListenAndServe(); err != nil {
				slog.Error(fmt.Errorf("smtp server: %w", err).Error())
			}
		}()
	}

	var submissionSrv *smtp.Server
	if config.Submission != nil {
		slog.Info("starting submission server", "address", config.Submission.Address)

		submissionSrv = &smtp.Server{
			Address:                 config.Submission.Address,
			Domain:                  "quokka.local",
			TLSConfig:               tlsConfig,
			AuthenticationEncrypted: true,
			AuthenticationMandatory: true,
		}

		go func() {
			if err := submissionSrv.ListenAndServe(); err != nil {
				slog.Error(fmt.Errorf("submission server: %w", err).Error())
			}
		}()
	}

	var submissionsSrv *smtp.Server
	if config.Submissions != nil {
		slog.Info("starting submissions server", "address", config.Submissions.Address)

		submissionsSrv = &smtp.Server{
			Address:                 config.Submissions.Address,
			Domain:                  "quokka.local",
			TLSConfig:               tlsConfig,
			AuthenticationEncrypted: true,
			AuthenticationMandatory: true,
		}

		go func() {
			if err := submissionsSrv.ListenAndServe(); err != nil {
				slog.Error(fmt.Errorf("submissions server: %w", err).Error())
			}
		}()
	}

	if config.IMAP != nil {
		slog.Info("starting imap server", "address", config.IMAP.Address)
	}

	if config.IMAPS != nil {
		slog.Info("starting imaps server", "address", config.IMAPS.Address)
	}

	var metricsServer *metrics.Server
	if config.Metrics != nil {
		slog.Info("starting metrics server", "address", config.Metrics.Address)

		metricsServer = &metrics.Server{
			Address: config.Metrics.Address,
		}

		go func() {
			if err := metricsServer.ListenAndServe(); err != nil {
				slog.Error(fmt.Errorf("metrics server: %w", err).Error())
			}
		}()
	}

	var pprofServer *pprof.Server
	if config.Pprof != nil {
		slog.Info("starting pprof server", "address", config.Pprof.Address)

		pprofServer = &pprof.Server{
			Address: config.Pprof.Address,
		}

		go func() {
			if err := pprofServer.ListenAndServe(); err != nil {
				slog.Error(fmt.Errorf("pprof server: %w", err).Error())
			}
		}()
	}

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)

	<-interruptSignal
	slog.Info("got an interrupt signal")
	ctx := context.Background()

	if smtpSrv != nil {
		if err := smtpSrv.Close(); err != nil {
			return err
		}
	}

	if submissionSrv != nil {
		if err := submissionSrv.Close(); err != nil {
			return err
		}
	}

	if submissionsSrv != nil {
		if err := submissionsSrv.Close(); err != nil {
			return err
		}
	}

	if metricsServer != nil {
		if err := metricsServer.Shutdown(ctx); err != nil {
			return err
		}
	}

	if pprofServer != nil {
		if err := pprofServer.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}
