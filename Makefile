# https://github.com/docker/docker-credential-helpers/blob/master/Makefile

all: test

test:
	go test -cover -v ./...

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o release/linux/amd64/dotfiles ./cmd/dotfiles

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o release/windows/dotfiles.exe ./cmd/dotfiles

darwin:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o release/darwin/dotfiles ./cmd/dotfiles
