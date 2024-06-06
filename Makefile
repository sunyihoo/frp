export PATH := $(PATH):`go env GOPATH`/bin
LDFLAGS := -s -w
print:
	echo $(PATH)

all: env fmt build

build: frps frpc

env:
	@go version

# compile assets into binary file
file:
	rm -rf  ./assets/frps/static/*
	rm -rf ./assets/frps/static/*
	cp -rf ./web/frps/dist/* ./assets/frps/static
	cp -rf ./web/frpc/dist/* ./assets/ftpc/static

fmt:
	go fmt ./...

fmt-more:
	gofumpt -l -w .

vet:
	go vet ./..

frps:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -tags frps -o bin/frps ./cmd/frps

frpc:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -tags frpc -o bin/fprc ./cmd/frpc

stat:
	git show --shortstat


clean:
	rm -f ./bin/frpc
	rm -f ./bin/frps
	rm -rf ./lastversion

