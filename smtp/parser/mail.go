package parser

import (
	"errors"
	"strings"
)

var (
	ErrMailCommandInvalid = errors.New("parser: Mail command is invalid")
)

type MailCommand struct {
	ReversePath string
}

func NewMailCommand(cmdAndArgs string) (*MailCommand, error) {
	if len(cmdAndArgs) < 10 || strings.ToUpper(cmdAndArgs[:10]) != "MAIL FROM:" {
		return nil, ErrMailCommandInvalid
	}

	args := strings.Split(strings.Trim(cmdAndArgs[10:], " "), " ")

	if args[0] == "" {
		return nil, ErrMailCommandInvalid
	}

	return &MailCommand{
		ReversePath: args[0],
	}, nil
}
