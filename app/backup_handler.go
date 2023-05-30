package app

import (
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type BackupHandler struct {
	*handler
}

func (h *BackupHandler) Run(args []string) error {
	h.handler = new(handler)
	err := h.init()
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

func (h *BackupHandler) backupEntry(filename string) error {
	log.Debug("backup: ", filename)
	entry := h.config.NewEntry(filename, false)
	if entry == nil {
		return nil
	}
	if entry.ignored {
		log.Debug("ignore: ", filename)
		return nil
	}
	if entry.isDir {
		return h.backupDir(entry)
	}
	return h.backupFile(entry)
}

func (h *BackupHandler) backupDir(entry *fileEntry) error {
	log.Debug("dir: ", entry.path)
	// read directory
	files, err := os.ReadDir(entry.homePath)
	if err != nil {
		return err
	}
	// loop files in the dir
	for _, file := range files {
		filename := path.Join(entry.path, file.Name())
		if err = h.backupEntry(filename); err != nil {
			return err
		}
	}
	return nil
}

func (h *BackupHandler) backupFile(entry *fileEntry) error {
	log.Debug("file: ", entry.path)
	// get file contents
	contents, err := os.ReadFile(entry.homePath)
	if err != nil {
		return err
	}

	f, err := openFile(entry.backupPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// backup file
	if entry.encrypted {
		log.Debug("encrypt: ", entry.path)
		contents, err = h.config.Encrypt.EncryptFile(f, contents)
		if err != nil {
			return err
		}
	}
	same, err := saveFile(f, contents)
	if err != nil {
		return err
	}
	if same {
		log.Debug("not modified: ", entry.path)
		return nil
	}
	log.Info("backup compeleted: ", entry.path)
	return nil
}

func (h *BackupHandler) pluginBackup(name string, conf pluginConfig) error {
	p, err := loadPlugin(name, conf)
	if err != nil {
		return err
	}
	contents, err := p.Backup()
	if err != nil {
		return err
	}
	filename := path.Join(h.config.Backup.Target, "plugins", name)
	if conf.Encrypt {
		filename += h.config.Encrypt.Extension()
	}
	f, err := openFile(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if conf.Encrypt {
		contents, err = h.config.Encrypt.EncryptFile(f, contents)
		if err != nil {
			return err
		}
	}
	same, err := saveFile(f, contents)
	if err != nil {
		return err
	}
	if same {
		log.Debug("not modified: ", filename)
		return nil
	}
	log.Info(name, " backup completed")
	return nil
}
