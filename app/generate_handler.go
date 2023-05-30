package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GenerateHandler struct {
	*Encrypt
}

func (h *GenerateHandler) Run(args []string) error {
	err := h.Encrypt.CreateKeyset()
	if err != nil {
		return err
	}
	log.Infof("%s created", h.Encrypt.KeyPath)
	return h.createYaml()
}

func (h *GenerateHandler) createYaml() error {
	example := &config{
		Name:    viper.GetString("name"),
		Version: viper.GetString("version"),
		Backup: &backupConfig{
			Target:  "./data/",
			Home:    []string{".bash_profile", ".ssh/"},
			Ignore:  []string{".ssh/authorized_keys", ".ssh/known_hosts"},
			Encrypt: []string{".ssh/*_rsa*"},
		},
		Encrypt: h.Encrypt,
	}
	data, err := json.Marshal(example)
	if err != nil {
		return err
	}
	return h.createConfig(viper.GetString("config"), bytes.NewReader(data))
}

func (h *GenerateHandler) createConfig(cfgFile string, in io.Reader) error {
	newViper := viper.New()
	newViper.SetConfigFile(cfgFile)
	if err := newViper.ReadConfig(in); err != nil {
		return err
	}
	// backup current config file
	if err := h.backupFile(cfgFile); err != nil {
		return err
	}
	// write config file
	if err := newViper.WriteConfig(); err != nil {
		return err
	}
	log.Infof("%s created", cfgFile)
	return nil
}

func (h *GenerateHandler) backupFile(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		// failed to load file means needn't backup, exit without error
		return nil
	}
	newpath := fmt.Sprintf("%s.%s", filename, time.Now().Format("20060102150405"))
	return os.Rename(filename, newpath)
}
