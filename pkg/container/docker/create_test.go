package docker

import (
	"archive/tar"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/UltimateSoftware/envctl/pkg/container"
)

func TestBuildDockerfile(t *testing.T) {
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
		prettyFail(t, "errors", "no errors", err)
	}

	expected := `FROM scratch
	VOLUME ["/test-path"]
	WORKDIR "/test-path"
	ENTRYPOINT ["/testsh"]`

	actual := buf.String()
	if expected != actual {
		prettyFail(t, "Dockerfile build", expected, actual)
	}
}

func TestGetBuildContext(t *testing.T) {
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
		prettyFail(t, "Dockerfile build", "no errors", err)
	}

	bldctx, err := getBuildContext(buf)
	if err != nil {
		prettyFail(t, "getBuildContext()", "no errors", err)
	}

	tarrd := tar.NewReader(bldctx)
	h, err := tarrd.Next()
	for err != io.EOF {
		// The only thing in here should be the Dockerfile. If this is the case,
		// the for loop will exit before hitting this check again. If it's not,
		// it will fail because whatever is next won't match this.
		if h.FileInfo().Name() != "Dockerfile" {
			prettyFail(t, "tar header", "Dockerfile", h.FileInfo().Name())
		}

		expected := int64(buf.Len())
		actual := h.FileInfo().Size()
		if expected != actual {
			prettyFail(t, "tar file size", expected, actual)
		}

		h, err = tarrd.Next()
	}
}

func prettyFail(t *testing.T, header string, expected, actual interface{}) {
	divider := strings.Repeat("=", len(header))

	output := fmt.Sprintf(
		"\n\n%v\n%v\n\nexpected:\n\n\t%v\n\ngot:\n\n\t%v\n",
		header,
		divider,
		expected,
		actual,
	)

	t.Fatal(output)
}
