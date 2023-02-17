package smtp_test

import (
	"reflect"
	"testing"

	"github.com/quokkamail/quokka/smtp"
)

func TestParseMailCommand(t *testing.T) {
	tests := []struct {
		name       string
		cmdAndArgs string
		want       *smtp.MailCommand
		wantErr    bool
	}{
		{
			name:       "InvalidName",
			cmdAndArgs: "MELA",
			wantErr:    true,
		},
		{
			name:       "NoArguments",
			cmdAndArgs: "MAIL",
			wantErr:    true,
		},
		{
			name:       "NoReversePath",
			cmdAndArgs: "MAIL FROM:",
			wantErr:    true,
		},
		{
			name:       "Valid",
			cmdAndArgs: "MAIL FROM: <foo@bar.com>",
			want: &smtp.MailCommand{
				ReversePath: "<foo@bar.com>",
			},
		},
		{
			name:       "ValidWithoutSpace",
			cmdAndArgs: "MAIL FROM:<foo@bar.com>",
			want: &smtp.MailCommand{
				ReversePath: "<foo@bar.com>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := smtp.ParseMailCommand(tt.cmdAndArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMailCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMailCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
