package dotfiles

import (
	"fmt"

	"github.com/fatindeed/dotfiles-go/mcrypt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type handler struct {
	config config
	vault  *vault
}

func (h *handler) init(c *cli.Context) error {
	// init vault
	cipher, err := mcrypt.NewCipher()
	if err != nil {
		return err
	}
	h.vault = &vault{
		cipher:  cipher,
		prefix:  fmt.Sprintf("%s-vault", c.App.Name),
		version: "1.0",
	}

	// get config file contents
	contents, err := fileGetContents(c.String("config"))
	if err != nil {
		return err
	}
	// parse config file
	if err = yaml.Unmarshal(contents, &h.config); err != nil {
		return err
	}
	if h.config.Name != c.App.Name {
		return fmt.Errorf("invalid app name (%s) in the config file", h.config.Name)
	}

	// init backup config
	return h.config.Backup.init()
}
