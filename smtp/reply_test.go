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

package smtp

import (
	"io"
	"net"
	"testing"

	"github.com/shoenig/test/must"
)

func Test_session_reply(t *testing.T) {
	t.Parallel()

	type args struct {
		code  uint
		lines []string
	}

	testCases := []struct {
		name string
		args args
		want string
	}{
		{
			name: "OneLine",
			args: args{
				code:  250,
				lines: []string{"Line1"},
			},
			want: "250 Line1\r\n",
		},
		{
			name: "MultipleLines",
			args: args{
				code:  250,
				lines: []string{"Line1", "Line2", "Line3"},
			},
			want: "250-Line1\r\n250-Line2\r\n250 Line3\r\n",
		},
		{
			name: "NoLines",
			args: args{
				code:  250,
				lines: []string{},
			},
			want: "250 \r\n",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, server := net.Pipe()

			s := &session{
				conn: server,
			}

			go func() {
				s.reply(tc.args.code, tc.args.lines...)

				server.Close()
			}()

			data, err := io.ReadAll(client)
			must.NoError(t, err)
			must.Eq(t, tc.want, string(data))
		})
	}
}
