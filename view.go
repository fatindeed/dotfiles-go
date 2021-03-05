package dotfiles

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatindeed/dotfiles-go/mcrypt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// ViewCommand returns the recover command
func ViewCommand() *cli.Command {
	return &cli.Command{
		Name:      "view",
		Aliases:   []string{"v"},
		Usage:     "view encrypted contents",
		ArgsUsage: "filename (- for stdin)",
		Before:    beforeAction(),
		Action: func(c *cli.Context) error {
			handler := &viewHandler{handler: &handler{}}
			return handler.run(c)
		},
		Flags: append(append(flags, mcrypt.CipherFlags()...), verboseFlag),
	}
}

type viewHandler struct {
	*handler
}

func (h *viewHandler) run(c *cli.Context) error {
	err := h.init(c)
	if err != nil {
		return err
	}
	log.Debug("view encrypted contents")
	args := c.Args()
	if args.Len() == 0 {
		return fmt.Errorf("missing filename")
	}
	filename := c.Args().Get(0)
	if filename == "-" {
		ciphertext, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		contents, err := h.vault.Decrypt(ciphertext)
		if err != nil {
			return err
		}
		fmt.Printf("%s", contents)
	} else {
		info, err := os.Stat(filename)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", filename)
		}
		contents, err := h.vault.DecryptFile(filename)
		if err != nil {
			return err
		}
		fmt.Printf("%s", contents)
	}
	return nil
}
