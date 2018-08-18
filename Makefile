GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test -v
GOCLEAN=$(GOCMD) clean
BINARY_NAME=envctl
VERSION=DEVEL
LDFLAGS=-X github.com/UltimateSoftware/envctl/cmd.envctlVersion=$(VERSION)
TARBALL_NAME=$(BINARY_NAME)$(VERSION).darwin-amd64.tar.gz

all: test $(TARBALL_NAME)

.PHONY: test
test:
	$(GOTEST) ./...

$(BINARY_NAME):
	$(GOBUILD) -v -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)

$(TARBALL_NAME): $(BINARY_NAME)
	tar -czvf $(TARBALL_NAME) $(BINARY_NAME)

.PHONY: clean
clean: $(BINARY_NAME) $(TARBALL_NAME)
	$(GOCLEAN)
	rm $(TARBALL_NAME)
