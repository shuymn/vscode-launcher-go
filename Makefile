ifdef update
	u=-u
endif

export GO111MODULE=on

all: build

build: clean
	go build -ldflags="-s -w"

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o vscode-launcher-go.exe

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
