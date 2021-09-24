BIN := vscode-launcher

ifdef update
	u=-u
endif

export GO111MODULE=on

all: build

build: $(BIN).go clean
	go build -ldflags="-s -w" $(BIN).go

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: test
test:
	go test -race ./...

.PHONY: clean
clean:
	@[ ! -f $(BIN) ] || rm $(BIN)