package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

// Encrypts specified plaintext with key.
// This implementation is based on example from https://golang.org/pkg/crypto/cipher/#example_NewCFBEncrypter
func Encrypt(plaintext []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	initVector := ciphertext[:aes.BlockSize]

	// Initialization vector needs to be unique but not secure, so fill it with random bytes.
	_, err = io.ReadFull(rand.Reader, initVector)
	if err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCFBEncrypter(block, initVector)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}

// Decrypts specified ciphertext with key.
// This implementation is based on example from https://golang.org/pkg/crypto/cipher/#example_NewCFBDecrypter
func Decrypt(ciphertext []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	if len(ciphertext) < aes.BlockSize {
		log.Fatal("Ciphertext is less than AES block size")
	}

	initVector := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCFBDecrypter(block, initVector)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext
}
