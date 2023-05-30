package app

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type ViewHandler struct {
	*handler
}

func (h *ViewHandler) Run(args []string) error {
	h.handler = new(handler)
	err := h.init()
	if err != nil {
		return err
	}
	log.Debug("view encrypted contents")
	if len(args) == 0 {
		return fmt.Errorf("missing filename")
	}
	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	contents, err := h.config.Encrypt.DecryptFile(f)
	if err != nil {
		return err
	}
	fmt.Printf("%s", contents)
	return nil
}
