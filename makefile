.PHONY: build

build:
	go tool vet .
	GOOS=darwin GOARCH=amd64 go build -o bin/dir github.com/adamdecaf/dist/dir
	chmod +x ./bin/*

test:
	go test -v ./dir
