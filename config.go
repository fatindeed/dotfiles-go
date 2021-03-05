package dotfiles

import (
	"fmt"
	"os"

	ignore "github.com/sabhiram/go-gitignore"
)

type config struct {
	Name    string
	Version string
	Backup  *backupConfig
}

type backupConfig struct {
	ignoreMatcher  *ignore.GitIgnore `yaml:"-"`
	encryptMatcher *ignore.GitIgnore `yaml:"-"`
	homeDir        string            `yaml:"-"`

	Target  string
	Home    []string
	Ignore  []string
	Encrypt []string
	Plugins map[string]pluginConfig
}

type pluginConfig struct {
	// simple plugin with given command
	Backup  []string
	Recover []string

	// custom plugin files with backup.ext and recover.ext
	Command   string
	Extension string

	// encrypt data
	Encrypt bool
}

func (c *backupConfig) init() error {
	// load home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	c.homeDir = fmt.Sprintf("%s/", homeDir)
	// init file rules
	c.ignoreMatcher = ignore.CompileIgnoreLines(c.Ignore...)
	c.encryptMatcher = ignore.CompileIgnoreLines(c.Encrypt...)
	return nil
}
