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

package parser_test

import (
	"reflect"
	"testing"

	"github.com/quokkamail/quokka/smtp/parser"
)

func TestNewMailCommand(t *testing.T) {
	type args struct {
		cmdAndArgs string
	}

	tests := []struct {
		name    string
		args    args
		want    *parser.MailCommand
		wantErr bool
	}{
		{
			name: "InvalidName",
			args: args{
				cmdAndArgs: "MELA",
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
				cmdAndArgs: "MAIL FROM:",
			},
			wantErr: true,
		},
		{
			name: "Valid",
			args: args{
				cmdAndArgs: "MAIL FROM: <foo@bar.com>",
			},
			want: &parser.MailCommand{
				ReversePath: "<foo@bar.com>",
			},
		},
		{
			name: "ValidWithoutSpace",
			args: args{
				cmdAndArgs: "MAIL FROM:<foo@bar.com>",
			},
			want: &parser.MailCommand{
				ReversePath: "<foo@bar.com>",
			},
		},
		{
			name: "ValidLowecase",
			args: args{
				cmdAndArgs: "mail from: <foo@bar.com>",
			},
			want: &parser.MailCommand{
				ReversePath: "<foo@bar.com>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.NewMailCommand(tt.args.cmdAndArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMailCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMailCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
