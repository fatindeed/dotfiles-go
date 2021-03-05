package dotfiles

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatindeed/dotfiles-go/mcrypt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// BackupCommand returns the backup command
func BackupCommand() *cli.Command {
	return &cli.Command{
		Name:    "backup",
		Aliases: []string{"b"},
		Usage:   "backup dotfiles",
		Before:  beforeAction(),
		Action: func(c *cli.Context) error {
			handler := &backupHandler{handler: &handler{}}
			return handler.run(c)
		},
		Flags: append(append(flags, mcrypt.CipherFlags()...), verboseFlag),
	}
}

type backupHandler struct {
	*handler
}

func (h *backupHandler) run(c *cli.Context) error {
	err := h.init(c)
	if err != nil {
		return err
	}
	log.Debug("backup files")
	for _, filename := range h.config.Backup.Home {
		err = h.backupEntry(filename)
		if err != nil {
			return err
		}
	}
	log.Debug("backup plugins")
	for name, conf := range h.config.Backup.Plugins {
		err = h.pluginBackup(name, conf)
		if err != nil {
			log.Error(err)
		}
	}
	log.Info("backup completed")
	return nil
}

func (h *backupHandler) backupEntry(filename string) error {
	log.Debug("backup: ", filename)
	entry := &fileEntry{config: h.config.Backup}
	entry.load(filename)
	src, err := fileStat(entry.homePath)
	if src == nil {
		return err
	}
	if src.IsDir() {
		return h.backupDir(entry)
	}
	if src.Mode().IsRegular() {
		return h.backupFile(entry)
	}
	return nil
}

func (h *backupHandler) backupDir(entry *fileEntry) error {
	if entry.ignored {
		log.Debug("ignore: ", entry.path)
		return nil
	}
	log.Debug("dir: ", entry.path)
	// read directory
	files, err := ioutil.ReadDir(entry.homePath)
	if err != nil {
		return err
	}
	if err = mkdir(entry.backupPath); err != nil {
		return err
	}
	// loop files in the dir
	dirname := strings.TrimRight(entry.path, "/")
	for _, file := range files {
		filename := fmt.Sprintf("%s/%s", dirname, file.Name())
		if err = h.backupEntry(filename); err != nil {
			return err
		}
	}
	return nil
}

func (h *backupHandler) backupFile(entry *fileEntry) error {
	if entry.ignored {
		log.Debug("ignore: ", entry.path)
		return nil
	}
	log.Debug("file: ", entry.path)
	if err := mkdir(filepath.Dir(entry.backupPath)); err != nil {
		return err
	}
	// backup file
	if entry.encrypted {
		return h.encryptFile(entry)
	}
	return copyFile(entry.homePath, entry.backupPath)
}

func (h *backupHandler) encryptFile(entry *fileEntry) error {
	log.Debug("encrypt: ", entry.path)
	// get file contents
	contents, err := fileGetContents(entry.homePath)
	if err != nil {
		return err
	}
	updated, err := h.vault.EncryptFile(entry.backupPath, contents)
	if err != nil {
		return err
	}
	if updated {
		log.Info("encrypted: ", entry.homePath, " => ", entry.backupPath)
	}
	return nil
}

func (h *backupHandler) pluginBackup(name string, conf pluginConfig) error {
	p, err := loadPlugin(name, conf)
	if err != nil {
		return err
	}
	contents, err := p.Backup()
	if err != nil {
		return err
	}
	pluginDir := path.Join(h.config.Backup.Target, "plugins")
	if err = mkdir(pluginDir); err != nil {
		return err
	}
	updated := true
	filename := path.Join(pluginDir, name)
	if conf.Encrypt {
		updated, err = h.vault.EncryptFile(filename, contents)
	} else {
		err = ioutil.WriteFile(filename, contents, 0644)
	}
	if err != nil {
		return err
	}
	if updated {
		log.Info(name, " backup completed")
	}
	return nil
}
