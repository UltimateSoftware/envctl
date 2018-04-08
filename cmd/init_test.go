package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/UltimateSoftware/envctl/test_pkg"
)

func TestInit(got *testing.T) {
	t := test_pkg.NewT(got)

	// This needs to be set in order for the init command to work. Normally it's
	// set by Cobra when the command is initialized.
	cfgFile = "envctl.yaml.test"

	cmd := newInitCmd()

	// Hijacking here swallows the command output so that it doesn't clutter
	// the output of `go test -v ./...`.
	outch, errch := test_pkg.HijackStdout(func() {
		cmd.Run(cmd, []string{})
	})

	select {
	case err := <-errch:
		t.Fatal("hijacking output", nil, err)
	case <-outch:
	}

	f, err := os.Open(cfgFile)
	if err != nil {
		t.Fatal("errors", nil, err)
	}

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal("errors", nil, err)
	}

	expected := `---
image: ubuntu:latest

shell: /bin/bash

bootstrap:
- echo 'Environment initialized' > /envctl

variables:
- FOO=bar
`

	actual := string(raw)

	if expected != actual {
		t.Fatal("file contents", expected, actual)
	}

	// This is a failure because the tests need to reset the state of the repo
	// back to how it was. Other tests rely on the repo being in a clean state.
	err = os.Remove(cfgFile)
	if err != nil {
		t.Fatal("error removing generated file", nil, err)
	}
}
