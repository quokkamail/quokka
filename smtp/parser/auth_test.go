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
