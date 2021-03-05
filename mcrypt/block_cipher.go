package mcrypt

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// blockCipher represents a block cipher running in a block-based mode (CBC, ECB etc).
type blockCipher struct {
	Block cipher.Block
	// Mode values:
	//
	// "cbc" (cipher block chaining) is a block cipher mode that is
	// significantly more secure than ECB mode.
	//
	// "ecb" (electronic codebook) is a block cipher mode that is generally
	// unsuitable for most purposes. The use of this mode is not recommended.
	// @see https://github.com/golang/go/issues/5597
	Mode string
}

func (c *blockCipher) newEncrypter(iv []byte) (mode cipher.BlockMode) {
	switch c.Mode {
	case "CBC":
		mode = cipher.NewCBCEncrypter(c.Block, iv)
	}
	return
}

func (c *blockCipher) newDecrypter(iv []byte) (mode cipher.BlockMode) {
	switch c.Mode {
	case "CBC":
		mode = cipher.NewCBCDecrypter(c.Block, iv)
	}
	return
}

// Encrypt encrypts the given byte slice and puts information about the final result in the returned value.
func (c *blockCipher) Encrypt(plaintext []byte) ([]byte, error) {
	blockSize := c.Block.BlockSize()

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	plaintext = padding(plaintext, blockSize)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, blockSize+len(plaintext))
	iv := ciphertext[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := c.newEncrypter(iv)
	if mode == nil {
		return nil, fmt.Errorf("block cipher init failed")
	}
	mode.CryptBlocks(ciphertext[blockSize:], plaintext)

	return ciphertext, nil
}

// Decrypt takes in the value and decrypts it into the byte slice.
func (c *blockCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	blockSize := c.Block.BlockSize()

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < blockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := c.newDecrypter(iv)
	if mode == nil {
		return nil, fmt.Errorf("block cipher init failed")
	}
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	return unpadding(ciphertext), nil
}
