package fs

import (
	"fmt"
)

// ReadFileFS is the interface implemented by a file system that provides an optimized implementation of ReadFile.
type ReadFileFS interface {
	ReadFile(name string) ([]byte, error)
}

func NewFilesystem(scheme string) (ReadFileFS, error) {
	switch scheme {
	case "op":
		return getOnepasswordAdapter()
	case "file":
		return getLocalFilesystemAdapter()
	}
	return nil, fmt.Errorf("%s not implemented", scheme)
}
