package cmd

import (
	"testing"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/UltimateSoftware/envctl/test_pkg"
)

func TestDestroy(got *testing.T) {
	t := test_pkg.NewT(got)

	cnt := container.Metadata{
		ID:        "foocnt",
		ImageID:   "fooimg",
		BaseName:  "fooenv",
		BaseImage: "scratch",
		Shell:     "/foo/sh",
		Mount: container.Mount{
			Destination: "/foo/mnt",
		},
	}

	s := &memStore{
		env: db.Environment{
			Status:    db.StatusReady,
			Container: cnt,
		},
	}

	ctl := newMockCtl(&cnt)

	cmd := newDestroyCmd(ctl, s)

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

	if nil != ctl.current {
		t.Fatal("backing container", nil, ctl.current)
	}

	if db.StatusOff != s.env.Status {
		t.Fatal("status", db.StatusOff, s.env.Status)
	}
}
