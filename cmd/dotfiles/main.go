package main

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/fatindeed/dotfiles-go/cmd"
	_ "github.com/joho/godotenv/autoload"
)

var (
	Version string
	Commit  string
)

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if Version == "" {
			Version = strings.Trim(info.Main.Version, "()")
		}
		if Commit == "" {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					Commit = setting.Value
					break
				}
			}
		}
	}
	shortCommit := ""
	if len(Commit) > 7 {
		shortCommit = Commit[0:7]
	}
	return strings.Trim(fmt.Sprintf("%s-%s", Version, shortCommit), "-")
}

func main() {
	cmd.Execute(getVersion())
}
