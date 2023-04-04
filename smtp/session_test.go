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

package smtp_test

import (
	"testing"
)

func TestCommandUnrecognized(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	if _, _, err := tc.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := tc.Cmd("DUMMYCOMMAND")
	if err != nil {
		t.Fatal(err)
	}

	tc.StartResponse(id)
	defer tc.EndResponse(id)

	if _, _, err := tc.ReadResponse(500); err != nil {
		t.Fatal(err)
	}
}

func TestHELOCommand(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	if _, _, err := tc.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := tc.Cmd("HELO")
	if err != nil {
		t.Fatal(err)
	}

	tc.StartResponse(id)
	defer tc.EndResponse(id)

	if _, _, err := tc.ReadResponse(250); err != nil {
		t.Fatal(err)
	}
}

func TestQUITCommand(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	if _, _, err := tc.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := tc.Cmd("QUIT")
	if err != nil {
		t.Fatal(err)
	}

	tc.StartResponse(id)
	defer tc.EndResponse(id)

	if _, _, err := tc.ReadResponse(221); err != nil {
		t.Fatal(err)
	}
}

func TestRSETCommand(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	if _, _, err := tc.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := tc.Cmd("RSET")
	if err != nil {
		t.Fatal(err)
	}

	tc.StartResponse(id)
	defer tc.EndResponse(id)

	if _, _, err := tc.ReadResponse(250); err != nil {
		t.Fatal(err)
	}
}

func TestNOOPCommand(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	if _, _, err := tc.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := tc.Cmd("NOOP")
	if err != nil {
		t.Fatal(err)
	}

	tc.StartResponse(id)
	defer tc.EndResponse(id)

	if _, _, err := tc.ReadResponse(250); err != nil {
		t.Fatal(err)
	}
}

func TestMAILCommandWithNoArguments(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	if _, _, err := tc.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := tc.Cmd("MAIL")
	if err != nil {
		t.Fatal(err)
	}

	tc.StartResponse(id)
	defer tc.EndResponse(id)

	if _, _, err := tc.ReadResponse(501); err != nil {
		t.Fatal(err)
	}
}

func TestMAILCommand(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	if _, _, err := tc.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := tc.Cmd("MAIL FROM: <mail@domain.ext>")
	if err != nil {
		t.Fatal(err)
	}

	tc.StartResponse(id)
	defer tc.EndResponse(id)

	if _, _, err := tc.ReadResponse(250); err != nil {
		t.Fatal(err)
	}
}
