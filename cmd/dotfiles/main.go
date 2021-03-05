package main

import (
	"os"

	"github.com/fatindeed/dotfiles-go"
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetOutput(os.Stdout)
	if err := gotenv.Load(); err != nil {
		log.Error("load env: ", err)
	}
	app := &cli.App{
		Name:    "dotfiles",
		Usage:   "A dotfiles management application",
		Version: "1.0.0",
		Flags:   []cli.Flag{},
		Commands: []*cli.Command{
			dotfiles.GenerateCommand(),
			dotfiles.BackupCommand(),
			dotfiles.RecoverCommand(),
			dotfiles.ViewCommand(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Error("runtime error: ", err)
	}
}
