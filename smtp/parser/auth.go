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
	ErrAuthCommandInvalid = errors.New("parser: Auth command is invalid")
)

type AuthCommand struct {
	Mechanism       string
	InitialResponse string
}

func NewAuthCommand(cmdAndArgs string) (*AuthCommand, error) {
	if len(cmdAndArgs) < 4 || strings.ToUpper(cmdAndArgs[:4]) != "AUTH" {
		return nil, ErrAuthCommandInvalid
	}

	args := strings.Split(strings.Trim(cmdAndArgs[4:], " "), " ")

	if args[0] == "" {
		return nil, ErrAuthCommandInvalid
	}

	mechanism := strings.ToUpper(args[0])
	initialResponse := ""
	if len(args) >= 2 {
		initialResponse = args[1]
	}

	return &AuthCommand{
		Mechanism:       mechanism,
		InitialResponse: initialResponse,
	}, nil
}
