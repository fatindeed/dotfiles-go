package dotfiles

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// GenerateCommand returns the generate command
func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"g"},
		Usage:   "Generate new config.yml and .env",
		Before:  beforeAction(),
		Action: func(c *cli.Context) error {
			handler := &generateHandler{}
			return handler.run(c)
		},
		Flags: append(flags, verboseFlag),
	}
}

type generateHandler struct{}

func (h *generateHandler) run(c *cli.Context) error {
	err := h.createYaml(c)
	if err != nil {
		return err
	}
	return h.createEnv()
}

func (h *generateHandler) createYaml(c *cli.Context) error {
	configPath := c.String("config")
	config := &config{
		Name:    c.App.Name,
		Version: c.App.Version,
		Backup: &backupConfig{
			Target:  "./data/",
			Home:    []string{".bash_profile", ".ssh/"},
			Ignore:  []string{".ssh/authorized_keys", ".ssh/known_hosts"},
			Encrypt: []string{".ssh/*_rsa*"},
		},
	}
	// yaml encode
	var data bytes.Buffer
	encoder := yaml.NewEncoder(&data)
	encoder.SetIndent(2)
	defer encoder.Close()
	if err := encoder.Encode(config); err != nil {
		return err
	}
	// backup current config file
	if err := h.backupFile(configPath); err != nil {
		return err
	}
	// file put contents
	if err := ioutil.WriteFile(configPath, data.Bytes(), 0644); err != nil {
		return err
	}
	log.Info("config.yml file created")
	return nil
}

func (h *generateHandler) createEnv() error {
	envPath := ".env"
	// backup current env file
	if err := h.backupFile(envPath); err != nil {
		return err
	}
	var data bytes.Buffer
	data.WriteString("DOTFILES_CIPHER_METHOD=AES-256-CBC\n")
	data.WriteString("DOTFILES_PASSPHRASE=" + randString(32) + "\n")
	// file put contents
	if err := ioutil.WriteFile(envPath, data.Bytes(), 0644); err != nil {
		return err
	}
	log.Info(".env file created")
	return nil
}

func (h *generateHandler) backupFile(filename string) error {
	if ok, err := fileExists(filename); !ok {
		return err
	}
	newpath := fmt.Sprintf("%s.%s", filename, time.Now().Format("20060102150405"))
	return os.Rename(filename, newpath)
}
