u := $(if $(update),-u)

all: build

build: clean
	CGO_ENABLED=0 go build -ldflags="-s -w"

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w"

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: clean
clean:
	@[ ! -f vscode-launcher-go ] || rm vscode-launcher-go
