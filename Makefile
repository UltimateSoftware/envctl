GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test -v
GOCLEAN=$(GOCMD) clean
BINARY_NAME=envctl

all: test $(BINARY_NAME)

.PHONY: test
test:
	$(GOTEST) ./...

$(BINARY_NAME):
	$(GOBUILD) -v -o $(BINARY_NAME)

.PHONY: clean
clean: $(BINARY_NAME)
	$(GOCLEAN)
