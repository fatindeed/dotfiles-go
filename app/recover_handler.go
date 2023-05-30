package app

import (
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type RecoverHandler struct {
	*handler
}

func (h *RecoverHandler) Run(args []string) error {
	h.handler = new(handler)
	err := h.init()
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

func (h *RecoverHandler) recoverEntry(filename string) error {
	log.Debug("recover: ", filename)
	entry := h.config.NewEntry(filename, true)
	if entry == nil {
		return nil
	}
	if entry.ignored {
		log.Debug("ignore: ", filename)
		return nil
	}
	if entry.isDir {
		return h.recoverDir(entry)
	}
	return h.recoverFile(entry)
}

func (h *RecoverHandler) recoverDir(entry *fileEntry) error {
	log.Debug("dir: ", entry.path)
	// read directory
	files, err := os.ReadDir(entry.backupPath)
	if err != nil {
		return err
	}
	// loop files in the dir
	for _, file := range files {
		filename := path.Join(entry.path, file.Name())
		if err = h.recoverEntry(filename); err != nil {
			return err
		}
	}
	return nil
}

func (h *RecoverHandler) recoverFile(entry *fileEntry) error {
	log.Debug("file: ", entry.path)
	var contents []byte
	if entry.encrypted {
		f, err := os.Open(entry.backupPath)
		if err != nil {
			return err
		}
		log.Debug("decrypt: ", entry.path)
		contents, err = h.config.Encrypt.DecryptFile(f)
		if err != nil {
			return err
		}
	} else {
		var err error
		contents, err = os.ReadFile(entry.backupPath)
		if err != nil {
			return err
		}
	}
	f, err := openFile(entry.homePath)
	if err != nil {
		return err
	}
	same, err := saveFile(f, contents)
	if err != nil {
		return err
	}
	if same {
		log.Debug("not modified: ", entry.path)
		return nil
	}
	log.Info("recover compeleted: ", entry.path)
	return nil
}

func (h *RecoverHandler) pluginRecover(name string, conf pluginConfig) error {
	p, err := loadPlugin(name, conf)
	if err != nil {
		return err
	}
	filename := path.Join(h.config.Backup.Target, "plugins", name)
	var contents []byte
	if conf.Encrypt {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		contents, err = h.config.Encrypt.DecryptFile(f)
		if err != nil {
			return err
		}
	} else {
		contents, err = os.ReadFile(filename)
		if err != nil {
			return err
		}
	}
	err = p.Recover(contents)
	if err != nil {
		return err
	}
	log.Info(name, " recover completed")
	return nil
}
