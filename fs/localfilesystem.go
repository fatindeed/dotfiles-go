package fs

import (
	"os"
)

type localFilesystemAdapter struct{}

func (c *localFilesystemAdapter) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

var localfsInstance *localFilesystemAdapter

func getLocalFilesystemAdapter() (*localFilesystemAdapter, error) {
	if localfsInstance != nil {
		return localfsInstance, nil
	}
	localfsInstance = new(localFilesystemAdapter)
	return localfsInstance, nil
}
