package smtp_test

import (
	"testing"
)

func TestCommandNotImplemented(t *testing.T) {
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

	if _, _, err := cst.c.ReadResponse(502); err != nil {
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
