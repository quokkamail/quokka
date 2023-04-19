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
	s.replyWithReply(replyReady(s.srv.Domain))

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
			s.replyWithReply(replyCommandUnrecognized())
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
		s.replyWithReply(replyAuthenticationRequired())
		return
	}

	if s.mailFrom != "" {
		s.replyWithReply(replyBadSequence())
		return
	}

	mailCmd, err := parser.NewMailCommand(cmdAndArgs)
	if err != nil {
		s.replyWithReply(replySyntaxError())
		return
	}

	s.mailFrom = mailCmd.ReversePath
	s.replyWithReply(replyOk())
}

func (s *session) handleRcptCommand(cmdAndArgs string) {
	if s.isNotAuthenticatedWhenMandatory() {
		s.replyWithReply(replyAuthenticationRequired())
		return
	}

	if s.mailFrom == "" {
		s.replyWithReply(replyBadSequence())
		return
	}

	recipientCmd, err := parser.NewRecipientCommand(cmdAndArgs)
	if err != nil {
		s.replyWithReply(replySyntaxError())
		return
	}

	s.rcptTo = append(s.rcptTo, recipientCmd.ForwardPath)
	s.replyWithReply(replyOk())
}

func (s *session) handleEhloCommand() {
	extensions := []string{"AUTH PLAIN", "PIPELINING"}
	if !s.tls {
		extensions = append(extensions, "STARTTLS")
	}

	s.replyWithReply(replyEhloOk(extensions))
}

func (s *session) handleHeloCommand() {
	s.replyWithReply(replyHeloOk())
}

func (s *session) handleQuitCommand() {
	s.replyWithReply(replyClosingConnection(s.srv.Domain))
	s.rwc.Close()
}

func (s *session) handleRsetCommand() {
	s.reset()
	s.replyWithReply(replyOk())
}

func (s *session) handleNoopCommand() {
	s.replyWithReply(replyOk())
}

func (s *session) handleDataCommand() {
	if s.isNotAuthenticatedWhenMandatory() {
		s.replyWithReply(replyAuthenticationRequired())
		return
	}

	if s.mailFrom == "" || len(s.rcptTo) == 0 {
		s.replyWithReply(replyBadSequence())
		return
	}

	s.replyWithReply(replyStartMailInput())

	dl, err := s.txtReader.ReadDotLines()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}

		slog.Error(fmt.Errorf("smtp session: %w", err).Error())
		return
	}

	s.data = dl
	s.replyWithReply(replyOk())
}

func (s *session) handleStartTLSCommand() {
	if s.tls {
		s.replyWithReply(replyBadSequence())
		return
	}

	s.replyWithReply(replyReadyToStartTLS())

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
		s.reply(503, "Bad sequence of commands")
		return
	}

	authCmd, err := parser.NewAuthCommand(cmdAndArgs)
	if err != nil {
		s.reply(501, "Syntax error in parameters or arguments")
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
		s.reply(235, "Authentication succeeded")
	default:
		s.reply(504, "Unrecognized authentication mechanism")
	}
}

func (s *session) isNotAuthenticatedWhenMandatory() bool {
	return s.srv.AuthenticationMandatory && !s.authenticated
}
