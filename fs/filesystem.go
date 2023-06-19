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
	case "s3":
		return getS3Adapter()
	case "op":
		return getOnepasswordAdapter()
	case "file":
		return getLocalStorageAdapter()
	}
	return nil, fmt.Errorf("%s not implemented", scheme)
}
