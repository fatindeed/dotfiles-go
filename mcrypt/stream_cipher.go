package mcrypt

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// streamCipher represents a stream cipher.
type streamCipher struct {
	Block cipher.Block
	// Mode values:
	//
	// "cfb" (cipher feedback, in 8-bit mode) is a stream cipher mode. It is
	// recommended to use NCFB mode rather than CFB mode.
	//
	// "ctr" (counter mode) is a stream cipher mode.
	//
	// "ofb" (output feedback, in 8-bit mode) is a stream cipher mode comparable
	// to CFB, but can be used in applications where error propagation cannot be
	// tolerated. It is recommended to use NOFB mode rather than OFB mode.
	Mode string
}

func (c *streamCipher) newEncrypter(iv []byte) (stream cipher.Stream) {
	switch c.Mode {
	case "CFB":
		stream = cipher.NewCFBEncrypter(c.Block, iv)
	case "CTR":
		stream = cipher.NewCTR(c.Block, iv)
	case "OFB":
		stream = cipher.NewOFB(c.Block, iv)
	}
	return
}

func (c *streamCipher) newDecrypter(iv []byte) (stream cipher.Stream) {
	switch c.Mode {
	case "CFB":
		stream = cipher.NewCFBDecrypter(c.Block, iv)
	case "CTR":
		stream = cipher.NewCTR(c.Block, iv)
	case "OFB":
		stream = cipher.NewOFB(c.Block, iv)
	}
	return
}

// Encrypt encrypts the given byte slice and puts information about the final result in the returned value.
func (c *streamCipher) Encrypt(plaintext []byte) ([]byte, error) {
	blockSize := c.Block.BlockSize()

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, blockSize+len(plaintext))
	iv := ciphertext[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := c.newEncrypter(iv)
	if stream == nil {
		return nil, fmt.Errorf("stream cipher init failed")
	}
	stream.XORKeyStream(ciphertext[blockSize:], plaintext)

	return ciphertext, nil
}

// Decrypt takes in the value and decrypts it into the byte slice.
func (c *streamCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	blockSize := c.Block.BlockSize()

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < blockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	stream := c.newDecrypter(iv)
	if stream == nil {
		return nil, fmt.Errorf("stream cipher init failed")
	}
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
