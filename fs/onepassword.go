package fs

import (
	"fmt"
	"strings"

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/1Password/connect-sdk-go/onepassword"
)

type onepasswordAdapter struct {
	connect.Client
	Vaults map[string]*onepasswordCache
}

type onepasswordCache struct {
	*onepassword.Vault
	Items map[string]*onepassword.Item
}

func (c *onepasswordAdapter) GetItem(vaultName, itemName string) (*onepassword.Item, error) {
	vaultCache, ok := c.Vaults[vaultName]
	if !ok {
		vault, err := c.GetVaultByTitle(vaultName)
		if err != nil {
			return nil, err
		}
		vaultCache = &onepasswordCache{
			Vault: vault,
			Items: map[string]*onepassword.Item{},
		}
		c.Vaults[vaultName] = vaultCache
	}

	item, ok := vaultCache.Items[itemName]
	if !ok {
		var err error
		item, err = c.GetItemByTitle(itemName, vaultCache.ID)
		if err != nil {
			return nil, err
		}
		vaultCache.Items[itemName] = item
	}
	return item, nil
}

// Read implements `op read`
func (c *onepasswordAdapter) ReadFile(name string) ([]byte, error) {
	parts := strings.Split(name, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid url: %s", name)
	}
	item, err := c.GetItem(parts[0], parts[1])
	if err != nil {
		return nil, err
	}
	for _, file := range item.Files {
		if file.Name != parts[2] {
			continue
		}
		b, err := c.GetFileContent(file)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	for _, field := range item.Fields {
		if field.Label != parts[2] {
			continue
		}
		return []byte(field.Value), nil
	}
	return nil, fmt.Errorf("%s not found", name)
}

var opInstance *onepasswordAdapter

func getOnepasswordAdapter() (*onepasswordAdapter, error) {
	if opInstance != nil {
		return opInstance, nil
	}
	c, err := connect.NewClientFromEnvironment()
	if err != nil {
		return nil, err
	}
	opInstance = &onepasswordAdapter{
		Client: c,
		Vaults: map[string]*onepasswordCache{},
	}
	return opInstance, nil
}
