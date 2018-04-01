package docker

import (
	"archive/tar"
	"io"
	"testing"

	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/UltimateSoftware/envctl/test_pkg"
)

func TestBuildDockerfile(got *testing.T) {
	t := test_pkg.NewT(got)

	testm := container.Metadata{
		BaseImage: "scratch",
		Mount: container.Mount{
			Source:      "",
			Destination: "/test-path",
		},
		Shell: "/testsh",
	}

	buf, err := buildDockerfile(testm)
	if err != nil {
		t.Fatal("errors", nil, err)
	}

	expected := `FROM scratch
	VOLUME ["/test-path"]
	WORKDIR "/test-path"
	ENTRYPOINT ["/testsh"]`

	actual := buf.String()
	if expected != actual {
		t.Fatal("Dockerfile build", expected, actual)
	}
}

func TestGetBuildContext(got *testing.T) {
	t := test_pkg.NewT(got)

	testm := container.Metadata{
		BaseImage: "scratch",
		Mount: container.Mount{
			Source:      "",
			Destination: "/test-path",
		},
		Shell: "/testsh",
	}

	buf, err := buildDockerfile(testm)
	if err != nil {
		t.Fatal("Dockerfile build", "no errors", err)
	}

	bldctx, err := getBuildContext(buf)
	if err != nil {
		t.Fatal("getBuildContext()", "no errors", err)
	}

	tarrd := tar.NewReader(bldctx)
	h, err := tarrd.Next()
	for err != io.EOF {
		// The only thing in here should be the Dockerfile. If this is the case,
		// the for loop will exit before hitting this check again. If it's not,
		// it will fail because whatever is next won't match this.
		if h.FileInfo().Name() != "Dockerfile" {
			t.Fatal("tar header", "Dockerfile", h.FileInfo().Name())
		}

		expected := int64(buf.Len())
		actual := h.FileInfo().Size()
		if expected != actual {
			t.Fatal("tar file size", expected, actual)
		}

		h, err = tarrd.Next()
	}
}
