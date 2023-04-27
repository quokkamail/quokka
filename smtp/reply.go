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

import "fmt"

const (
	genericOkReply              string = "2.0.0 Requested mail action okay, completed"
	badSequenceReply            string = "5.5.1 Bad sequence of commands"
	syntaxErrorReply            string = "5.5.4 Syntax error in parameters or arguments"
	authenticationRequiredReply string = "5.7.0 Authentication required"
)

func (s *session) reply(code uint, lines ...string) {
	if len(lines) == 0 {
		fmt.Fprintf(s.rwc, "%d \r\n", code)
	}

	for _, m := range lines[:len(lines)-1] {
		fmt.Fprintf(s.rwc, "%d-%s\r\n", code, m)
	}
	fmt.Fprintf(s.rwc, "%d %s\r\n", code, lines[len(lines)-1])
}
