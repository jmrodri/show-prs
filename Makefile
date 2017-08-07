.DEFAULT_GOAL    := build
SOURCES          := $(shell find . -name '*.go' -not -path "*/vendor/*")

show-prs: $(SOURCES)
	go build -i -ldflags="-s -w"

build: show-prs
	@echo > /dev/null

clean:
	@rm -f show-prs

test:
	go test ./pkg/...

.PHONY: build clean test
