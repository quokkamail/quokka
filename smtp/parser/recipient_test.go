package parser_test

import (
	"reflect"
	"testing"

	"github.com/quokkamail/quokka/smtp/parser"
)

func TestNewRecipientCommand(t *testing.T) {
	type args struct {
		cmdAndArgs string
	}

	tests := []struct {
		name    string
		args    args
		want    *parser.RecipientCommand
		wantErr bool
	}{
		{
			name: "InvalidName",
			args: args{
				cmdAndArgs: "RECI",
			},
			wantErr: true,
		},
		{
			name: "NoArguments",
			args: args{
				cmdAndArgs: "MAIL",
			},
			wantErr: true,
		},
		{
			name: "NoReversePath",
			args: args{
				cmdAndArgs: "RCPT TO:",
			},
			wantErr: true,
		},
		{
			name: "Valid",
			args: args{
				cmdAndArgs: "RCPT TO: <foo@bar.com>",
			},
			want: &parser.RecipientCommand{
				ForwardPath: "<foo@bar.com>",
			},
		},
		{
			name: "ValidWithoutSpace",
			args: args{
				cmdAndArgs: "RCPT TO:<foo@bar.com>",
			},
			want: &parser.RecipientCommand{
				ForwardPath: "<foo@bar.com>",
			},
		},
		{
			name: "ValidLowecase",
			args: args{
				cmdAndArgs: "rcpt to: <foo@bar.com>",
			},
			want: &parser.RecipientCommand{
				ForwardPath: "<foo@bar.com>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.NewRecipientCommand(tt.args.cmdAndArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRecipientCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRecipientCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
