package util

import (
	"testing"
)

func TestEncryptionDeryption(t *testing.T) {
	key := []byte("This is a secret")
	plaintext := []byte("Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.")
	ciphertext := Encrypt(plaintext, key)

	if AssertArraysAreEqual(plaintext, ciphertext) {
		t.Fatal("Plaintext and ciphertext are the same!")
	}

	decryptedCiphertext := Decrypt(ciphertext, key)
	if !AssertArraysAreEqual(decryptedCiphertext, plaintext) {
		t.Fatal("Decrypted ciphertext does not match plaintext!")
	}

	wrongKey := []byte("This key should NOT work")
	decryptedCiphertext = Decrypt(ciphertext, wrongKey)
	if AssertArraysAreEqual(decryptedCiphertext, plaintext) {
		t.Fatal("Ciphertext was decrypted with wrong key!")
	}
}
