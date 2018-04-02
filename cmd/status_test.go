package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/test_pkg"
)

type memStore struct {
	env db.Environment
}

func (s *memStore) Create(e db.Environment) error {
	s.env = e

	return nil
}

func (s *memStore) Read() (db.Environment, error) {
	return s.env, nil
}

func (s *memStore) Delete() error {
	s.env = db.Environment{}
	return nil
}

func TestOffStatus(got *testing.T) {
	t := test_pkg.NewT(got)

	oldstdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("initializing test", nil, err)
	}

	os.Stdout = w

	s := &memStore{
		env: db.Environment{
			Status: 0,
		},
	}

	cmd := newStatusCmd(s)

	cmd.Run(cmd, []string{})

	os.Stdout = oldstdout

	errch := make(chan error)
	outch := make(chan []byte)
	go func() {
		hijackedStdout, err := ioutil.ReadAll(r)
		if err != nil {
			errch <- err
		}

		outch <- hijackedStdout
		r.Close()
	}()

	w.Close()

	expected := `The environment is off.

Run "envctl create" to spin it up!
`

	select {
	case err := <-errch:
		t.Fatal("error reading output", nil, err)
	case actual := <-outch:
		if expected != string(actual) {
			t.Fatal("output", expected, string(actual))
		}
	}
}
