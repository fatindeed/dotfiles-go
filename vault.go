package dotfiles

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fatindeed/dotfiles-go/mcrypt"
)

// https://golang.org/src/crypto/cipher/example_test.go
// https://github.com/ansible/ansible/blob/devel/lib/ansible/parsing/vault/__init__.py
// https://segmentfault.com/a/1190000021267253

type vault struct {
	cipher mcrypt.Cipher

	prefix  string
	version string
}

// Encrypt data
func (v *vault) Encrypt(plaintext []byte) ([]byte, error) {
	if v.isEncrypted(plaintext) {
		return nil, fmt.Errorf("input is already encrypted")
	}

	ciphertext, err := v.cipher.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return v.pack(ciphertext)
}

// Decrypt data
func (v *vault) Decrypt(data []byte) ([]byte, error) {
	if !v.isEncrypted(data) {
		return nil, fmt.Errorf("input is not vault encrypted data")
	}

	ciphertext, err := v.unpack(data)
	if err != nil {
		return nil, err
	}

	return v.cipher.Decrypt(ciphertext)
}

// EncryptFile encrypt a file
func (v *vault) EncryptFile(filename string, plaintext []byte) (bool, error) {
	// decrypt file contents
	contents, err := v.DecryptFile(filename)
	// skip if not modified
	if err == nil && bytes.Equal(plaintext, contents) {
		return false, nil
	}
	// encrypt contents
	ciphertext, err := v.Encrypt(plaintext)
	if err != nil {
		return false, err
	}
	// file put contents
	return true, ioutil.WriteFile(filename, ciphertext, 0644)
}

// DecryptFile decrypt a file
func (v *vault) DecryptFile(filename string) ([]byte, error) {
	// get file contents
	ciphertext, err := fileGetContents(filename)
	if err != nil {
		return nil, err
	}
	// decrypt contents
	return v.Decrypt(ciphertext)
}

// Test if this is vault encrypted data blob
func (v *vault) isEncrypted(data []byte) bool {
	slices := strings.Split(string(data), ";")
	if len(slices) != 3 || slices[0] != v.prefix {
		return false
	}
	return true
}

func (v *vault) isEncryptedFile(name string) (bool, error) {
	data, err := fileGets(name, 30)
	if err != nil {
		return false, err
	}
	return v.isEncrypted(data), nil
}

func (v *vault) pack(ciphertext []byte) ([]byte, error) {
	parts := []string{v.prefix, v.version,
		base64.StdEncoding.EncodeToString(ciphertext)}

	var payload bytes.Buffer
	_, err := payload.WriteString(strings.Join(parts, ";"))
	if err != nil {
		return nil, err
	}

	return payload.Bytes(), nil
}

func (v *vault) unpack(contents []byte) ([]byte, error) {
	slices := strings.Split(string(contents), ";")
	if len(slices) != 3 || slices[0] != v.prefix {
		return nil, fmt.Errorf("unable to parse envelope")
	}
	data, err := base64.StdEncoding.DecodeString(slices[2])
	return data, err
}
