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
		return nil, ErrMailCommandInvalid
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
