# https://github.com/docker/docker-credential-helpers/blob/master/Makefile

all: test

test:
	go test -cover -v ./...

linuxrelease:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o release/linux/amd64/dotfiles cmd/dotfiles/main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o release/linux/arm64/dotfiles cmd/dotfiles/main.go
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o release/linux/arm/dotfiles cmd/dotfiles/main.go

winrelease:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o release/windows/dotfiles.exe cmd/dotfiles/main.go

osxrelease:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o release/darwin/dotfiles cmd/dotfiles/main.go
