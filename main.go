package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	rsaEncryptor := NewRSAEncryptor(2048)

	// Save the public and private keys to files
	rsaEncryptor.SaveKeysToFile("public.pem", "private.pem")

	secretMessage := "Halo halo halo"

	// Encrypt the message
	encryptedMessage := rsaEncryptor.Encrypt(secretMessage)
	fmt.Println("Cipher Text:", encryptedMessage)

	// Decrypt the message
	rsaEncryptor.Decrypt(encryptedMessage)

	// Implementasi secrecy
	serverPublicKey, err := loadPublicKey("public.pem")
	CheckError(err)

	// Klien kirim pesan ke server, berisikan session key yang sudah dienkrip dengan server_public.key.
	sessionKey := "thisisaverysecretkey"
	encryptedSessionKey := rsaEncryptor.Encrypt(sessionKey)
	fmt.Println("\nKlien mengirim pesan ke server:")
	fmt.Println("Encrypted Session Key:", encryptedSessionKey)

	// Server decode
	decryptedSessionKey := rsaEncryptor.Decrypt(encryptedSessionKey)
	fmt.Println("\nServer mendecode pesan:")
	fmt.Println("Decrypted Session Key:", decryptedSessionKey)
	fmt.Println("Server jawab \"ok1\"")

	// Klien menerima "ok1".
	fmt.Println("\nKlien menerima \"ok1\".")

	// Klien kirim pesan ke server, data yang sudah dienkrip dengan session key.
	serverEncryptedMessage := RSAEncryptWithKey(secretMessage, *serverPublicKey)
	fmt.Println("\nKlien mengirim pesan ke server:")
	fmt.Println("Encrypted Message:", serverEncryptedMessage)

	// Server terima pesan
	serverDecryptedMessage := rsaEncryptor.Decrypt(serverEncryptedMessage)
	fmt.Println("\nServer mendecode pesan:")
	fmt.Println("Decrypted Message:", serverDecryptedMessage)
	fmt.Println("Server jawab \"ok2\"")

	// Klien menerima "ok2", cetak di layar
	fmt.Println("\nKlien menerima \"ok2\".")
}

// RSAEncryptWithKey encrypts a message using RSA-OAEP with a given public key.
func RSAEncryptWithKey(message string, publicKey rsa.PublicKey) string {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &publicKey, []byte(message), label)
	CheckError(err)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// loadPublicKey loads a public key from a PEM file.
func loadPublicKey(filename string) (*rsa.PublicKey, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(file)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of public key to *rsa.PublicKey")
	}

	return publicKey, nil
}
