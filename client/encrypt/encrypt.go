package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func hashPassword(pw string) []byte {
	hash := sha256.Sum256([]byte(pw))
	return hash[:]
}

func Encrypt(data []byte, password string) ([]byte, error) {
	block, err := aes.NewCipher(hashPassword(password))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	cipher := gcm.Seal(nonce, nonce, data, nil)
	return cipher, nil
}

func Decrypt(data []byte, password string) ([]byte, error) {
	block, err := aes.NewCipher(hashPassword(password))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func EncryptToFile(filename string, data []byte, password string) error {
	cipher, err := Encrypt(data, password)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(cipher)
	if err != nil {
		return err
	}
	return f.Sync()
}

func DecryptFromFile(filename string, password string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Decrypt(data, password)
}

func DecryptByStdin(filename string) ([]byte, error) {
	fmt.Printf("Password of %s: ", filename)
	pw, err := terminal.ReadPassword(0)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n")
	return DecryptFromFile(filename, string(pw))
}
