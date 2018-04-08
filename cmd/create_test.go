package cmd

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/UltimateSoftware/envctl/internal/db"
	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/UltimateSoftware/envctl/test_pkg"
)

func TestCreate(got *testing.T) {
	t := test_pkg.NewT(got)

	s := &memStore{
		env: db.Environment{
			Status: db.StatusOff,
		},
	}

	cfg := viper.New()

	cfg.Set("image", "test")
	cfg.Set("shell", "/foo/sh")
	cfg.Set("mount", "/foo/mnt")

	ctl := newMockCtl()

	cmd := newCreateCmd(ctl, s, cfg)

	cmd.Run(cmd, []string{})

	expected := db.Environment{
		Status: db.StatusReady,
		Container: container.Metadata{
			BaseImage: "test",
			Shell:     "/foo/sh",
			Mount: container.Mount{
				Destination: "/foo/mnt",
			},
		},
	}

	actual := s.env

	if expected.Status != actual.Status {
		t.Fatal("environment status", expected.Status, actual.Status)
	}

	if expected.Container.BaseImage != actual.Container.BaseImage {
		t.Fatal("environment image",
			expected.Container.BaseImage, actual.Container.BaseImage)
	}

	if expected.Container.Shell != actual.Container.Shell {
		t.Fatal("environment shell",
			expected.Container.Shell, actual.Container.Shell)
	}

	if expected.Container.Mount.Destination !=
		actual.Container.Mount.Destination {

		t.Fatal(
			"environment mount point",
			expected.Container.Mount.Destination,
			actual.Container.Mount.Destination,
		)
	}
}
