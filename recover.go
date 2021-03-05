package dotfiles

import (
	"bytes"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatindeed/dotfiles-go/mcrypt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// RecoverCommand returns the recover command
func RecoverCommand() *cli.Command {
	return &cli.Command{
		Name:    "recover",
		Aliases: []string{"r"},
		Usage:   "recover dotfiles",
		Before:  beforeAction(),
		Action: func(c *cli.Context) error {
			handler := &recoverHandler{handler: &handler{}}
			return handler.run(c)
		},
		Flags: append(append(flags, mcrypt.CipherFlags()...), verboseFlag),
	}
}

type recoverHandler struct {
	*handler
}

func (h *recoverHandler) run(c *cli.Context) error {
	err := h.init(c)
	if err != nil {
		return err
	}
	log.Debug("recover files")
	for _, filename := range h.config.Backup.Home {
		err = h.recoverEntry(filename)
		if err != nil {
			return err
		}
	}
	log.Debug("recover plugins")
	for name, conf := range h.config.Backup.Plugins {
		err = h.pluginRecover(name, conf)
		if err != nil {
			log.Error(err)
		}
	}
	log.Info("recover completed")
	return nil
}

func (h *recoverHandler) recoverEntry(filename string) error {
	log.Debug("recover: ", filename)
	entry := &fileEntry{config: h.config.Backup}
	entry.load(filename)
	src, err := fileStat(entry.backupPath)
	if src == nil {
		return err
	}
	if src.IsDir() {
		return h.recoverDir(entry)
	}
	if src.Mode().IsRegular() {
		return h.recoverFile(entry)
	}
	return nil
}

func (h *recoverHandler) recoverDir(entry *fileEntry) error {
	if entry.ignored {
		log.Debug("ignore: ", entry.path)
		return nil
	}
	log.Debug("dir: ", entry.path)
	// read directory
	files, err := ioutil.ReadDir(entry.backupPath)
	if err != nil {
		return err
	}
	// loop files in the dir
	dirname := strings.TrimRight(entry.path, "/")
	for _, file := range files {
		// TODO: trim .enc correctly
		filename := path.Join(dirname, strings.TrimRight(file.Name(), encExtension))
		if err = h.recoverEntry(filename); err != nil {
			return err
		}
	}
	return nil
}

func (h *recoverHandler) recoverFile(entry *fileEntry) error {
	if entry.ignored {
		log.Debug("ignore: ", entry.path)
		return nil
	}
	log.Debug("file: ", entry.path)
	if err := mkdir(filepath.Dir(entry.homePath)); err != nil {
		return err
	}
	// recover file
	if entry.encrypted {
		isEncrypted, err := h.vault.isEncryptedFile(entry.backupPath)
		if err != nil {
			return err
		}
		if isEncrypted {
			return h.decryptFile(entry)
		}
	}
	return copyFile(entry.backupPath, entry.homePath)
}

func (h *recoverHandler) decryptFile(entry *fileEntry) error {
	log.Debug("decrypt: ", entry.path)
	// decrypt contents
	data, err := h.vault.DecryptFile(entry.backupPath)
	if err != nil {
		return err
	}
	// check if home file changed
	contents, err := fileGetContents(entry.homePath)
	// skip if not modified
	if err == nil && bytes.Equal(contents, data) {
		log.Debug("not modified: ", entry.path)
		return nil
	}
	// file put contents
	err = ioutil.WriteFile(entry.homePath, data, 0644)
	if err == nil {
		log.Info("decrypted: ", entry.backupPath, " => ", entry.homePath)
	}
	return err
}

func (h *recoverHandler) pluginRecover(name string, conf pluginConfig) error {
	p, err := loadPlugin(name, conf)
	if err != nil {
		return err
	}
	filename := path.Join(h.config.Backup.Target, "plugins", name)
	var data []byte
	if conf.Encrypt {
		data, err = h.vault.DecryptFile(filename)
	} else {
		data, err = fileGetContents(filename)
	}
	if err != nil {
		return err
	}
	err = p.Recover(data)
	if err != nil {
		return err
	}
	log.Info(name, " recover completed")
	return nil
}
