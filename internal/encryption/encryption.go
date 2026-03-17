package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"filippo.io/age"
	"golang.org/x/crypto/argon2"
)

func ArgonHash(password []byte) ([]byte, []byte) {
	hash := rand.Text()
	key := argon2.IDKey(password, []byte(hash), 1, 64*1024, 4, 32)
	return key, []byte(hash)
}

func EncryptAge(data []byte, masterKey string) ([]byte, error) {
	recipient, err := age.NewScryptRecipient(masterKey)
	if err != nil {
		return nil, ErrCreateAgeRecipient(err)
	}

	out := &bytes.Buffer{}

	writer, err := age.Encrypt(out, recipient)
	if err != nil {
		return nil, err
	}

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func DecryptAge(masterKey string, encryptedData []byte) ([]byte, error) {
	identity, err := age.NewScryptIdentity(masterKey)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(encryptedData)

	r, err := age.Decrypt(reader, identity)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func EncryptGCM(key, name string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, []byte(name), nil), nil
}

func DecryptGCM(key, ciphertext string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
}
