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
	cmd.Run(cmd, []string{})

	if nil != ctl.current {
		t.Fatal("backing container", nil, ctl.current)
	}

	if db.StatusOff != s.env.Status {
		t.Fatal("status", db.StatusOff, s.env.Status)
	}
}
