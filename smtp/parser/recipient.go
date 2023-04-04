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

package parser

import (
	"errors"
	"strings"
)

var (
	ErrRecipientCommandInvalid = errors.New("parser: Recipient command is invalid")
)

type RecipientCommand struct {
	ForwardPath string
}

func NewRecipientCommand(cmdAndArgs string) (*RecipientCommand, error) {
	if len(cmdAndArgs) < 8 || strings.ToUpper(cmdAndArgs[:8]) != "RCPT TO:" {
		return nil, ErrRecipientCommandInvalid
	}

	args := strings.Split(strings.Trim(cmdAndArgs[8:], " "), " ")

	if args[0] == "" {
		return nil, ErrRecipientCommandInvalid
	}

	return &RecipientCommand{
		ForwardPath: args[0],
	}, nil
}
