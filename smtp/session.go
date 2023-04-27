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

package smtp

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strings"

	"github.com/quokkamail/quokka/smtp/parser"
	"golang.org/x/exp/slog"
)

type session struct {
	authenticated bool
	rwc           net.Conn
	srv           *Server
	tls           bool
	txtReader     *textproto.Reader

	data     []string
	mailFrom string
	rcptTo   []string
}

func (s *session) serve() {
	s.reply(220, fmt.Sprintf("%s ESMTP service ready", s.srv.Domain))

	s.txtReader = textproto.NewReader(bufio.NewReader(s.rwc))

	for {
		// s.rwc.SetReadDeadline()

		cmdAndArgs, err := s.txtReader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			slog.Error(fmt.Errorf("smtp session: %w", err).Error())
			return
		}

		command, _, _ := strings.Cut(cmdAndArgs, " ")

		switch strings.ToUpper(command) {
		case "EHLO":
			s.handleEhloCommand()
		case "HELO":
			s.handleHeloCommand()
		case "MAIL":
			s.handleMailCommand(cmdAndArgs)
		case "RCPT":
			s.handleRcptCommand(cmdAndArgs)
		case "DATA":
			s.handleDataCommand()
		case "QUIT":
			s.handleQuitCommand()
		case "RSET":
			s.handleRsetCommand()
		case "NOOP":
			s.handleNoopCommand()
		case "STARTTLS":
			s.handleStartTLSCommand()
		case "AUTH":
			s.handleAuthCommand(cmdAndArgs)
		default:
			s.reply(500, "5.5.2 Syntax error, command unrecognized")
		}
	}
}

func (s *session) reset() {
	s.rcptTo = make([]string, 0)
	s.mailFrom = ""
	s.data = make([]string, 0)
}

func (s *session) handleMailCommand(cmdAndArgs string) {
	if s.isNotAuthenticatedWhenMandatory() {
		s.reply(530, authenticationRequiredReply)
		return
	}

	if s.mailFrom != "" {
		s.reply(503, badSequenceReply)
		return
	}

	mailCmd, err := parser.NewMailCommand(cmdAndArgs)
	if err != nil {
		s.reply(501, syntaxErrorReply)
		return
	}

	s.mailFrom = mailCmd.ReversePath
	s.reply(250, "2.1.0 Requested mail action okay, completed")
}

func (s *session) handleRcptCommand(cmdAndArgs string) {
	if s.isNotAuthenticatedWhenMandatory() {
		s.reply(530, authenticationRequiredReply)
		return
	}

	if s.mailFrom == "" {
		s.reply(503, badSequenceReply)
		return
	}

	recipientCmd, err := parser.NewRecipientCommand(cmdAndArgs)
	if err != nil {
		s.reply(501, syntaxErrorReply)
		return
	}

	s.rcptTo = append(s.rcptTo, recipientCmd.ForwardPath)
	s.reply(250, "2.1.5 Requested mail action okay, completed")
}

func (s *session) handleEhloCommand() {
	extensions := []string{
		"Hello, nice to meet you",
		"AUTH PLAIN", "ENHANCEDSTATUSCODES", "PIPELINING",
	}

	if !s.tls {
		extensions = append(extensions, "STARTTLS")
	}

	s.reply(250, extensions...)
}

func (s *session) handleHeloCommand() {
	s.reply(250, "Hello, nice to meet you")
}

func (s *session) handleQuitCommand() {
	s.reply(221, fmt.Sprintf("2.0.0 %s service closing transmission channel", s.srv.Domain))
	s.rwc.Close()
}

func (s *session) handleRsetCommand() {
	s.reset()
	s.reply(250, genericOkReply)
}

func (s *session) handleNoopCommand() {
	s.reply(250, genericOkReply)
}

func (s *session) handleDataCommand() {
	if s.isNotAuthenticatedWhenMandatory() {
		s.replyWithReply(replyAuthenticationRequired())
		return
	}

	if s.mailFrom == "" || len(s.rcptTo) == 0 {
		s.reply(503, badSequenceReply)
		return
	}

	s.reply(354, "Start mail input; end with <CRLF>.<CRLF>")

	dl, err := s.txtReader.ReadDotLines()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}

		slog.Error(fmt.Errorf("smtp session: %w", err).Error())
		return
	}

	s.data = dl
	s.reply(250, genericOkReply)
}

func (s *session) handleStartTLSCommand() {
	if s.tls {
		s.reply(503, badSequenceReply)
		return
	}

	s.reply(220, "Ready to start TLS")

	tlsConn := tls.Server(s.rwc, s.srv.TLSConfig)
	if err := tlsConn.Handshake(); err != nil {
		slog.Error(fmt.Errorf("smtp session: %w", err).Error())
		s.replyWithReply(replyTLSNotAvailable())
		return
	}

	s.rwc = tlsConn
	s.txtReader = textproto.NewReader(bufio.NewReader(s.rwc))
	s.tls = true
	s.reset()
}

func (s *session) handleAuthCommand(cmdAndArgs string) {
	if s.srv.AuthenticationEncrypted && !s.tls {
		s.reply(530, "Must issue a STARTTLS command first")
		return
	}

	if s.authenticated {
		s.reply(503, badSequenceReply)
		return
	}

	authCmd, err := parser.NewAuthCommand(cmdAndArgs)
	if err != nil {
		s.reply(501, syntaxErrorReply)
		return
	}

	initialResponseBytes, err := base64.StdEncoding.DecodeString(authCmd.InitialResponse)
	if err != nil {
		s.reply(501, "Cannot decode response")
		slog.Error(fmt.Errorf("smtp session: %w", err).Error())
		return
	}
	initialResponse := string(initialResponseBytes)

	switch authCmd.Mechanism {
	case "PLAIN":
		if initialResponse == "" {
			s.reply(334)

			initialResponse, err = s.txtReader.ReadLine()
			if err != nil {
				slog.Error(fmt.Errorf("smtp session: %w", err).Error())
				return
			}
		}

		initialResponseParts := strings.Split(initialResponse, string([]byte{0}))
		if len(initialResponseParts) < 3 {
			return
		}

		// TODO: real authentication
		s.reply(235, "2.7.0 Authentication succeeded")
	default:
		s.reply(504, "Unrecognized authentication mechanism")
	}
}

func (s *session) isNotAuthenticatedWhenMandatory() bool {
	return s.srv.AuthenticationMandatory && !s.authenticated
}
