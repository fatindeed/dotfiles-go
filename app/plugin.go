package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	log "github.com/sirupsen/logrus"
)

type plugin interface {
	Backup() ([]byte, error)
	Recover([]byte) error
}

func loadPlugin(name string, conf pluginConfig) (plugin, error) {
	simple := true
	if len(conf.Backup) == 0 || len(conf.Recover) == 0 {
		pluginDir := path.Join("plugins", name)
		info, err := os.Stat(pluginDir)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("%s is not a dir", pluginDir)
		}
		simple = false
	}
	return &pluginHandler{
		name:   name,
		simple: simple,
		config: conf,
	}, nil
}

type pluginHandler struct {
	name   string
	simple bool
	config pluginConfig
}

func (p *pluginHandler) getCommand(action string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	if p.simple {
		var command []string
		switch action {
		case "backup":
			command = p.config.Backup
		case "recover":
			command = p.config.Recover
		default:
			return nil, fmt.Errorf("invalid plugin action: %s", action)
		}
		name, err := exec.LookPath(command[0])
		if err != nil {
			return nil, err
		}
		cmd = exec.Command(name, command[1:]...)
	} else {
		// load plugin command
		name, err := exec.LookPath(p.config.Command)
		if err != nil {
			return nil, err
		}
		// load plugin filename
		if p.config.Extension == "" {
			p.config.Extension = p.config.Command
		}
		filename := action + "." + p.config.Extension
		pluginFile := path.Join("plugins", p.name, filename)
		cmd = exec.Command(name, pluginFile)
	}
	return cmd, nil
}

// Backup with plugin command
func (p *pluginHandler) Backup() ([]byte, error) {
	cmd, err := p.getCommand("backup")
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Error(out.String())
	}
	return out.Bytes(), err
}

// Recover with plugin command
func (p *pluginHandler) Recover(data []byte) error {
	cmd, err := p.getCommand("recover")
	if err != nil {
		return err
	}
	var out bytes.Buffer
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Error(out.String())
	}
	return err
}
