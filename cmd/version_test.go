package cmd

import (
	"regexp"
	"testing"

	"github.com/UltimateSoftware/envctl/test_pkg"
)

func TestVersion(got *testing.T) {
	t := test_pkg.NewT(got)
	envctlVersion = "foo"

	cmd := newVersionCmd()

	outch, errch := test_pkg.HijackStdout(func() {
		cmd.Run(cmd, []string{})
	})

	expected := "foo\n"

	select {
	case err := <-errch:
		t.Fatal("hijacking output", nil, err)
	case actual := <-outch:
		if expected != string(actual) {
			t.Fatal("output", expected, string(actual))
		}
	}
}

func TestMissingVersion(got *testing.T) {
	t := test_pkg.NewT(got)

	cmd := newVersionCmd()

	outch, errch := test_pkg.HijackStdout(func() {
		cmd.Run(cmd, []string{})
	})

	match := "no version set for this build"

	select {
	case err := <-errch:
		t.Fatal("hijacking output", nil, err)
	case output := <-outch:
		if match, err := regexp.Match(match, output); !match || err != nil {
			t.Fatal("output contents", match, string(output))
		}
	}
}
