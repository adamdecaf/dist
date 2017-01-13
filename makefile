.PHONY: build

build:
	go tool vet .
	GOOS=darwin GOARCH=amd64 go build -o bin/dir github.com/adamdecaf/dist/dir
# 	Workers
	GOOS=darwin GOARCH=amd64 go build -o bin/worker-math github.com/adamdecaf/dist/workers/math
	chmod +x ./bin/*

test:
	go test -v ./dir
