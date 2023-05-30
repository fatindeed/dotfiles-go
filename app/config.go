package app

import (
	"fmt"
	"os"
	"path"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

type config struct {
	Name    string
	Version string
	Backup  *backupConfig
	Encrypt *Encrypt
	Context *configContext `mapstructure:"-"`
}

type configContext struct {
	homeDir        string
	ignoreMatcher  *ignore.GitIgnore
	encryptMatcher *ignore.GitIgnore
}

type backupConfig struct {
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

type fileEntry struct {
	path       string
	isDir      bool
	ignored    bool
	encrypted  bool
	homePath   string
	backupPath string
}

func (c *config) init() error {
	c.Context = new(configContext)
	// load home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	c.Context.homeDir = fmt.Sprintf("%s/", homeDir)
	// init file rules
	c.Context.ignoreMatcher = ignore.CompileIgnoreLines(c.Backup.Ignore...)
	c.Context.encryptMatcher = ignore.CompileIgnoreLines(c.Backup.Encrypt...)
	// init encryption engine
	return c.Encrypt.Init(c.Backup.Target)
}

func (c *config) NewEntry(filename string, isRecover bool) *fileEntry {
	entry := &fileEntry{path: filename}
	if isRecover {
		entry.path = strings.TrimSuffix(entry.path, c.Encrypt.Extension())
	}
	// load home entry
	entry.homePath = path.Join(c.Context.homeDir, entry.path)
	if c.Context.ignoreMatcher.MatchesPath(entry.path) {
		entry.ignored = true
		return entry
	}

	statfile := entry.homePath
	if isRecover {
		statfile = path.Join(c.Backup.Target, filename)
	}
	stat, err := os.Stat(statfile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return nil
	}
	if stat.IsDir() {
		entry.isDir = true
	}

	if c.Context.encryptMatcher.MatchesPath(entry.path) {
		entry.encrypted = true
	}
	// load backup entry
	entry.backupPath = path.Join(c.Backup.Target, entry.path)
	if !entry.isDir && entry.encrypted {
		entry.backupPath += c.Encrypt.Extension()
	}
	return entry
}
