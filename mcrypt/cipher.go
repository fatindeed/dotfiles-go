package mcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	algorithm  string
	passphrase string
)

// CipherFlags return cli flags
func CipherFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "algorithm",
			Aliases:     []string{"a"},
			Usage:       "the cipher method",
			Destination: &algorithm,
			EnvVars:     []string{"DOTFILES_ALGORITHM"},
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "passphrase",
			Aliases:     []string{"p"},
			Usage:       "the passphrase",
			Destination: &passphrase,
			EnvVars:     []string{"DOTFILES_PASSPHRASE"},
			Required:    true,
		},
	}
}

// Cipher is the embedded implementation to encrypting and decrypting data
type Cipher interface {
	// Encrypt encrypts the given byte slice and puts information about the final result in the returned value.
	Encrypt([]byte) ([]byte, error)
	// Decrypt takes in the value and decrypts it into the byte slice.
	Decrypt([]byte) ([]byte, error)
}

func newCipherBlock(ciphername string) (cipher.Block, error) {
	key := []byte(passphrase)
	switch ciphername {
	case "AES-128", "AES-192", "AES-256":
		keyLens := map[string]int{
			"AES-128": 16,
			"AES-192": 24,
			"AES-256": 32,
		}
		if keyLen := keyLens[ciphername]; len(key) != keyLen {
			return nil, fmt.Errorf("key must be %d bytes for %s", keyLen, algorithm)
		}
		return aes.NewCipher(key)
	case "DES":
		return des.NewCipher(key)
	}
	return nil, fmt.Errorf("invalid cipher method: %s", algorithm)
}

// NewCipher returns a new cipher
func NewCipher() (Cipher, error) {
	algorithm = strings.ToUpper(algorithm)
	slices := strings.Split(algorithm, "-")
	last := len(slices) - 1
	// init cipher block
	ciphername := strings.Join(slices[:last], "-")
	block, err := newCipherBlock(ciphername)
	if err != nil {
		return nil, err
	}
	// init block mode
	mode := slices[last]
	var cipher Cipher
	switch mode {
	case "CBC":
		cipher = &blockCipher{Block: block, Mode: mode}
	case "CFB", "OFB", "CTR":
		cipher = &streamCipher{Block: block, Mode: mode}
	case "GCM":
		cipher = &gcmCipher{Block: block}
	}
	if cipher == nil {
		return nil, fmt.Errorf("invalid cipher method: %s", algorithm)
	}
	return cipher, nil
}
