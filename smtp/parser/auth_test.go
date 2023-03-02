package parser_test

import (
	"reflect"
	"testing"

	"github.com/quokkamail/quokka/smtp/parser"
)

func TestNewAuthCommand(t *testing.T) {
	type args struct {
		cmdAndArgs string
	}

	tests := []struct {
		name    string
		args    args
		want    *parser.AuthCommand
		wantErr bool
	}{
		{
			name: "InvalidName",
			args: args{
				cmdAndArgs: "AUTO",
			},
			wantErr: true,
		},
		{
			name: "NoArguments",
			args: args{
				cmdAndArgs: "AUTH",
			},
			wantErr: true,
		},
		{
			name: "Valid",
			args: args{
				cmdAndArgs: "AUTH PLAIN",
			},
			want: &parser.AuthCommand{
				Mechanism: "PLAIN",
			},
		},
		{
			name: "ValidWithInitialResponse",
			args: args{
				cmdAndArgs: "AUTH PLAIN =",
			},
			want: &parser.AuthCommand{
				Mechanism:       "PLAIN",
				InitialResponse: "=",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.NewAuthCommand(tt.args.cmdAndArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAuthCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
