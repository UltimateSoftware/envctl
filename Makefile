GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test -v
BINARY_NAME=envctl

test:
	$(GOTEST) ./...
