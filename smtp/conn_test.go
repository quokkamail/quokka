package smtp_test

import (
	"testing"
)

func TestCommandUnrecognized(t *testing.T) {
	t.Parallel()

	cst := newClientServerTest(t)

	if _, _, err := cst.c.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := cst.c.Cmd("DUMMYCOMMAND")
	if err != nil {
		t.Fatal(err)
	}

	cst.c.StartResponse(id)
	defer cst.c.EndResponse(id)

	if _, _, err := cst.c.ReadResponse(500); err != nil {
		t.Fatal(err)
	}
}

func TestHELOCommand(t *testing.T) {
	t.Parallel()

	cst := newClientServerTest(t)

	if _, _, err := cst.c.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := cst.c.Cmd("HELO")
	if err != nil {
		t.Fatal(err)
	}

	cst.c.StartResponse(id)
	defer cst.c.EndResponse(id)

	if _, _, err := cst.c.ReadResponse(250); err != nil {
		t.Fatal(err)
	}
}

func TestQUITCommand(t *testing.T) {
	t.Parallel()

	cst := newClientServerTest(t)

	if _, _, err := cst.c.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := cst.c.Cmd("QUIT")
	if err != nil {
		t.Fatal(err)
	}

	cst.c.StartResponse(id)
	defer cst.c.EndResponse(id)

	if _, _, err := cst.c.ReadResponse(221); err != nil {
		t.Fatal(err)
	}
}

func TestRSETCommand(t *testing.T) {
	t.Parallel()

	cst := newClientServerTest(t)

	if _, _, err := cst.c.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := cst.c.Cmd("RSET")
	if err != nil {
		t.Fatal(err)
	}

	cst.c.StartResponse(id)
	defer cst.c.EndResponse(id)

	if _, _, err := cst.c.ReadResponse(250); err != nil {
		t.Fatal(err)
	}
}

func TestNOOPCommand(t *testing.T) {
	t.Parallel()

	cst := newClientServerTest(t)

	if _, _, err := cst.c.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := cst.c.Cmd("NOOP")
	if err != nil {
		t.Fatal(err)
	}

	cst.c.StartResponse(id)
	defer cst.c.EndResponse(id)

	if _, _, err := cst.c.ReadResponse(250); err != nil {
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

func TestMAILCommandWithInvalidArguments(t *testing.T) {
	t.Parallel()

	ts := NewTestServer(t)
	tc := ts.Client()

	if _, _, err := tc.ReadResponse(220); err != nil {
		t.Fatal(err)
	}

	id, err := tc.Cmd("MAIL ARG1")
	if err != nil {
		t.Fatal(err)
	}

	tc.StartResponse(id)
	defer tc.EndResponse(id)

	if _, _, err := tc.ReadResponse(501); err != nil {
		t.Fatal(err)
	}
}
