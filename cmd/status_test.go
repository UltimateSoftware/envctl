package cmd

import (
	"testing"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/test_pkg"
)

func TestOffStatus(got *testing.T) {
	t := test_pkg.NewT(got)

	s := &memStore{
		env: db.Environment{
			Status: db.StatusOff,
		},
	}

	cmd := newStatusCmd(s)

	outch, errch := test_pkg.HijackStdout(func() {
		cmd.Run(cmd, []string{})
	})

	expected := `The environment is off.

Run "envctl create" to spin it up!
`

	select {
	case err := <-errch:
		t.Fatal("hijacking output", nil, err)
	case actual := <-outch:
		if expected != string(actual) {
			t.Fatal("output", expected, string(actual))
		}
	}
}

func TestReadyStatus(got *testing.T) {
	t := test_pkg.NewT(got)
	s := &memStore{
		env: db.Environment{
			Status: db.StatusReady,
		},
	}

	cmd := newStatusCmd(s)

	outch, errch := test_pkg.HijackStdout(func() {
		cmd.Run(cmd, []string{})
	})

	expected := `The environment is ready!

Run "envctl login" to enter it.
`

	select {
	case err := <-errch:
		t.Fatal("hijacking output", nil, err)
	case actual := <-outch:
		if expected != string(actual) {
			t.Fatal("output", expected, string(actual))
		}
	}
}

func TestErrorStatus(got *testing.T) {
	t := test_pkg.NewT(got)
	s := &memStore{
		env: db.Environment{
			Status: db.StatusError,
		},
	}

	cmd := newStatusCmd(s)

	outch, errch := test_pkg.HijackStdout(func() {
		cmd.Run(cmd, []string{})
	})

	expected := `Something is wrong with the environment. :(

Try recreating it by running "envctl destroy", followed by "envctl create".
`

	select {
	case err := <-errch:
		t.Fatal("hijacking output", nil, err)
	case actual := <-outch:
		if expected != string(actual) {
			t.Fatal("output", expected, string(actual))
		}
	}
}
