package docker

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuildDockerfile(t *testing.T) {
	testcfg := ImageConfig{
		BaseImage: "scratch",
		Mount:     "/test-path",
		Shell:     "/testsh",
	}

	buf, err := buildDockerfile(testcfg)
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
