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
