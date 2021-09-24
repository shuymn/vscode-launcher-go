ifdef update
	u=-u
endif

export GO111MODULE=on

all: build

build: vscode-launcher_darwin.go clean
	go build -ldflags="-s -w"

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: test
test:
	go test -race ./...

.PHONY: clean
clean:
	@[ ! -f vscode-launcher-go ] || rm vscode-launcher-go