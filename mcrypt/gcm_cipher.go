package mcrypt

import (
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// gcmCipher represents a Galois Counter Mode with a specific key. See
// https://csrc.nist.gov/groups/ST/toolkit/BCM/documents/proposedmodes/gcm/gcm-revised-spec.pdf
type gcmCipher struct {
	Block cipher.Block
}

// Encrypt encrypts the given byte slice and puts information about the final result in the returned value.
func (c *gcmCipher) Encrypt(plaintext []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(c.Block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt takes in the value and decrypts it into the byte slice.
func (c *gcmCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(c.Block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}
