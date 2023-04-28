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
	"testing"

	"github.com/quokkamail/quokka/smtp/parser"
	"github.com/shoenig/test/must"
)

func TestNewRecipientCommand(t *testing.T) {
	t.Parallel()

	type args struct {
		cmdAndArgs string
	}

	testCases := []struct {
		name    string
		args    args
		want    *parser.RecipientCommand
		wantErr error
	}{
		{
			name: "InvalidName",
			args: args{
				cmdAndArgs: "RECI",
			},
			wantErr: parser.ErrRecipientCommandInvalid,
		},
		{
			name: "NoArguments",
			args: args{
				cmdAndArgs: "MAIL",
			},
			wantErr: parser.ErrRecipientCommandInvalid,
		},
		{
			name: "NoReversePath",
			args: args{
				cmdAndArgs: "RCPT TO:",
			},
			wantErr: parser.ErrRecipientCommandInvalid,
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

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := parser.NewRecipientCommand(tc.args.cmdAndArgs)
			must.Eq(t, tc.wantErr, err)
			must.Eq(t, tc.want, got)
		})
	}
}
