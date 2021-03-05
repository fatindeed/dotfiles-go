package dotfiles

import (
	"fmt"
	"strings"
)

const encExtension = ".enc"

type fileEntry struct {
	config *backupConfig

	path       string
	ignored    bool
	encrypted  bool
	homePath   string
	backupPath string
}

func (f *fileEntry) load(filename string) {
	f.path = filename
	// load home file
	f.homePath = fmt.Sprintf("%s%s", f.config.homeDir, f.path)
	// load backup file
	if f.config.ignoreMatcher.MatchesPath(f.path) {
		f.ignored = true
		return
	}
	if f.config.encryptMatcher.MatchesPath(f.path) {
		f.encrypted = true
	}
	f.backupPath = fmt.Sprintf("%s%s", f.config.Target, f.path)
	if !strings.HasSuffix(f.path, "/") && f.encrypted {
		f.backupPath += encExtension
	}
}
