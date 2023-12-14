GO_BUILD=CGO_ENABLED=0 go build -buildvcs

all: test release-build
	gzip -f build/*

test:
	go test crypto/

release-build: linux-arm64 linux-amd64

linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -o build/secret-server-$@

linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO_BUILD) -o build/secret-server-$@

clean:
	rm -rf build
