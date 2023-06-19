package fs

import (
	"os"
)

type localStorageAdapter struct{}

func (c *localStorageAdapter) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

var localfsInstance *localStorageAdapter

func getLocalStorageAdapter() (*localStorageAdapter, error) {
	if localfsInstance != nil {
		return localfsInstance, nil
	}
	localfsInstance = new(localStorageAdapter)
	return localfsInstance, nil
}
