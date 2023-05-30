package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/insecurecleartextkeyset"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/proto/tink_go_proto"
	"github.com/google/tink/go/tink"
)

// Encrypt and decrypt files and data.
type Encrypt struct {
	tink.AEAD `mapstructure:"-"`
	baseDir   string `mapstructure:"-"`
	// Configurable fields
	Template     string
	KeyPath      string
	Passphrase   string
	MasterKeyURI string `mapstructure:",omitempty"`
}

func (e *Encrypt) CreateKeyset() error {
	// https://github.com/google/tink/blob/master/testing/go/keyset_service.go
	var kt *tink_go_proto.KeyTemplate
	switch e.Template {
	case "AES128_GCM":
		kt = aead.AES128GCMKeyTemplate()
	case "AES256_GCM":
		kt = aead.AES256GCMKeyTemplate()
	case "AES256_GCM_RAW":
		kt = aead.AES256GCMNoPrefixKeyTemplate()
	// case "AES128_GCM_SIV":
	// 	kt = aead.AES128GCMSIVKeyTemplate()
	// case "AES256_GCM_SIV":
	// 	kt = aead.AES256GCMSIVKeyTemplate()
	// case "AES256_GCM_SIV_RAW":
	// 	kt = aead.AES256GCMSIVNoPrefixKeyTemplate()
	case "AES128_CTR_HMAC_SHA256":
		kt = aead.AES128CTRHMACSHA256KeyTemplate()
	case "AES256_CTR_HMAC_SHA256":
		kt = aead.AES256CTRHMACSHA256KeyTemplate()
	case "CHACHA20_POLY1305":
		kt = aead.ChaCha20Poly1305KeyTemplate()
	case "XCHACHA20_POLY1305":
		kt = aead.XChaCha20Poly1305KeyTemplate()
	default:
		return fmt.Errorf("unknown key template: %s", e.Template)
	}
	kh, err := keyset.NewHandle(kt)
	if err != nil {
		return err
	}

	f, err := os.Create(e.KeyPath)
	if err != nil {
		return err
	}
	defer f.Close()

	jw := keyset.NewJSONWriter(f)
	return insecurecleartextkeyset.Write(kh, jw)
}

func (e *Encrypt) Init(baseDir string) error {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return err
	}
	e.baseDir = fmt.Sprintf("%s/", strings.TrimRight(absBaseDir, "/"))

	b, err := getFileContents(e.KeyPath)
	if err != nil {
		return err
	}
	r := bytes.NewBuffer(b)
	kh, err := insecurecleartextkeyset.Read(keyset.NewJSONReader(r))
	if err != nil {
		return err
	}
	e.AEAD, err = aead.New(kh)
	return err
}

func (e *Encrypt) EncryptFile(f *os.File, plaintext []byte) ([]byte, error) {
	contents, err := e.DecryptFile(f)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(contents, plaintext) {
		return io.ReadAll(f)
	}
	return e.Encrypt(plaintext, []byte(e.Passphrase))
}

func (e *Encrypt) DecryptFile(f *os.File) ([]byte, error) {
	ciphertext, err := io.ReadAll(f)
	if err != nil || len(ciphertext) == 0 {
		return ciphertext, err
	}
	return e.Decrypt(ciphertext, []byte(e.Passphrase))
}

func (e *Encrypt) Extension() string {
	return ".aead"
}
