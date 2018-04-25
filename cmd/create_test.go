package cmd

import (
	"os"
	"testing"

	"github.com/UltimateSoftware/envctl/internal/config"
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

	cfg := memConfig{
		opts: config.Opts{
			Image: "test",
			Shell: "/foo/sh",
			Mount: "/foo/mnt",
		},
	}

	ctl := newMockCtl(nil)

	cmd := newCreateCmd(ctl, s, cfg)

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

	expectedStatus := db.StatusReady
	expectedContainer := container.Metadata{
		BaseImage: "test",
		Shell:     "/foo/sh",
		Mount: container.Mount{
			Destination: "/foo/mnt",
		},
	}

	// Testing that the user-specified configuration is saved correctly.
	if expectedStatus != s.env.Status {
		t.Fatal("environment status", expectedStatus, s.env.Status)
	}

	if expectedContainer.BaseImage != s.env.Container.BaseImage {
		t.Fatal("environment image",
			expectedContainer.BaseImage, s.env.Container.BaseImage)
	}

	if expectedContainer.Shell != s.env.Container.Shell {
		t.Fatal("environment shell",
			expectedContainer.Shell, s.env.Container.Shell)
	}

	if expectedContainer.Mount.Destination !=
		s.env.Container.Mount.Destination {

		t.Fatal(
			"environment mount point",
			expectedContainer.Mount.Destination,
			s.env.Container.Mount.Destination,
		)
	}

	// Now that correct saving of user-specified configuration has been
	// established, the calls to the container engine can be tested to make
	// sure that what's done there is totally in sync with what's been saved.
	if s.env.Container.ID != ctl.current.ID {
		t.Fatal("container id", s.env.Container.ID, ctl.current.ID)
	}

	if s.env.Container.ImageID != ctl.current.ImageID {
		t.Fatal(
			"container image id",
			s.env.Container.ImageID,
			ctl.current.ImageID,
		)
	}

	if s.env.Container.BaseImage != ctl.current.BaseImage {
		t.Fatal(
			"container base image",
			s.env.Container.BaseImage,
			ctl.current.BaseImage,
		)
	}

	if s.env.Container.BaseName != ctl.current.BaseName {
		t.Fatal("container base name",
			s.env.Container.BaseName,
			ctl.current.BaseName,
		)
	}

	if s.env.Container.Shell != ctl.current.Shell {
		t.Fatal("container shell",
			s.env.Container.Shell,
			ctl.current.Shell,
		)
	}

	if s.env.Container.Mount.Destination != ctl.current.Mount.Destination {
		t.Fatal("container mount point",
			s.env.Container.Mount.Destination,
			ctl.current.Mount.Destination,
		)
	}
}

// TODO: implement this
// func TestCreateDefaultMount(got *testing.T) {}

// TODO: implement this
// func TestCreateAlreadyInitialized(got *testing.T) {}

func TestCreateWithVariables(got *testing.T) {
	t := test_pkg.NewT(got)

	s := &memStore{
		env: db.Environment{
			Status: db.StatusOff,
		},
	}

	cfg := memConfig{
		opts: config.Opts{
			Image: "test",
			Shell: "/foo/sh",
			Mount: "/foo/mnt",
			Variables: map[string]string{
				"foo": "bar",
			},
		},
	}

	ctl := newMockCtl(nil)

	cmd := newCreateCmd(ctl, s, cfg)

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

	expected := "foo=bar"
	if s.env.Container.Envs[0] != expected {
		t.Fatal("variables", expected, s.env.Container.Envs[0])
	}
}

func TestCreateWithDynamicVariables(got *testing.T) {
	t := test_pkg.NewT(got)

	s := &memStore{
		env: db.Environment{
			Status: db.StatusOff,
		},
	}

	cfg := memConfig{
		opts: config.Opts{
			Image: "test",
			Shell: "/foo/sh",
			Mount: "/foo/mnt",
			Variables: map[string]string{
				"ENVCTL_TESTING": "$ENVCTL_TESTING",
			},
		},
	}

	ctl := newMockCtl(nil)

	os.Setenv("ENVCTL_TESTING", "FOO")
	defer os.Setenv("ENVCTL_TESTING", "")

	cmd := newCreateCmd(ctl, s, cfg)

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

	expected := "ENVCTL_TESTING=FOO"
	if s.env.Container.Envs[0] != expected {
		t.Fatal("variables", expected, s.env.Container.Envs[0])
	}
}

func TestParseMissingVariables(got *testing.T) {
	t := test_pkg.NewT(got)

	opts := config.Opts{
		Image: "test",
		Shell: "/foo/sh",
		Mount: "/foo/mnt",
		Variables: map[string]string{
			"ENVCTL_TESTING": "$ENVCTL_TESTING",
		},
	}

	envs, err := parseVariables(opts)
	if err == nil {
		t.Fatal("error parsing variables", "missing variable ENVCTL_TESTING", err)
	}

	if len(envs) != 0 {
		t.Fatal("number of parsed missing variables", 0, len(envs))
	}
}

func TestNoCache(got *testing.T) {
	t := test_pkg.NewT(got)

	cfg := memConfig{
		opts: config.Opts{
			Image:      "test",
			Shell:      "/foo/sh",
			Mount:      "/foo/mnt",
			CacheImage: config.NoCacheImage,
		},
	}

	ctl := newMockCtl(nil)

	s := &memStore{
		env: db.Environment{
			Status: db.StatusOff,
		},
	}

	cmd := newCreateCmd(ctl, s, cfg)

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

	if s.env.Container.NoCache != true {
		t.Fatal("setting nocache", true, s.env.Container.NoCache)
	}
}
