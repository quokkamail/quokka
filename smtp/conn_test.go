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
