package dbcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/xerrors"
)

type Cipher interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type DecryptFailedError struct {
	Inner error
}

func (e *DecryptFailedError) Error() string {
	return xerrors.Errorf("decrypt failed: %w", e.Inner).Error()
}

func (e *DecryptFailedError) Unwrap() error {
	return e.Inner
}

func IsDecryptFailedError(err error) bool {
	var e *DecryptFailedError
	return errors.As(err, &e)
}

// CipherAES256 returns a new AES-256 cipher.
func CipherAES256(key []byte) (Cipher, error) {
	if len(key) != 32 {
		return nil, xerrors.Errorf("key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &aes256{aead}, nil
}

type aes256 struct {
	aead cipher.AEAD
}

func (a *aes256) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, a.aead.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	return a.aead.Seal(nonce, nonce, plaintext, nil), nil
}

func (a *aes256) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < a.aead.NonceSize() {
		return nil, xerrors.Errorf("ciphertext too short")
	}
	decrypted, err := a.aead.Open(nil, ciphertext[:a.aead.NonceSize()], ciphertext[a.aead.NonceSize():], nil)
	if err != nil {
		return nil, &DecryptFailedError{Inner: err}
	}
	return decrypted, nil
}
