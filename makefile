.PHONY: build

build:
	go tool vet .
	cd dir && GOOS=darwin GOARCH=386 go build -o ../bin/dir .

test:
	cd dir && go test -v .
