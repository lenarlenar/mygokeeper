package storage

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	passphrase := "test-secret"
	original := []byte("this is a test")

	encrypted, err := encrypt(original, passphrase)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	if bytes.Equal(encrypted, original) {
		t.Error("encrypted data should differ from original")
	}

	decrypted, err := decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, original) {
		t.Errorf("decrypted data does not match original\nGot:  %s\nWant: %s", decrypted, original)
	}
}

func TestDecryptWithWrongPassphrase(t *testing.T) {
	pass1 := "right-password"
	pass2 := "wrong-password"
	plaintext := []byte("super secret")

	ciphertext, err := encrypt(plaintext, pass1)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	_, err = decrypt(ciphertext, pass2)
	if err == nil {
		t.Fatal("expected decryption to fail with wrong passphrase")
	}
}
