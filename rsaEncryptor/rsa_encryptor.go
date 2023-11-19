package rsaEncryptor

import (
	"Tugas1/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

// RSAEncryptor is a struct that encapsulates RSA public and private keys.
type RSAEncryptor struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// NewRSAEncryptor generates a new RSAEncryptor with a new key pair.
func NewRSAEncryptor(bits int) *RSAEncryptor {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	utils.CheckError(err)

	return &RSAEncryptor{
		PublicKey:  &privateKey.PublicKey,
		PrivateKey: privateKey,
	}
}

// SaveKeysToFile saves the public and private keys to files.
func (r *RSAEncryptor) SaveKeysToFile(publicKeyFile, privateKeyFile string) {
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(r.PublicKey)
	utils.CheckError(err)
	r.saveKeyToFile(publicKeyFile, "PUBLIC KEY", publicKeyPEM)
	r.saveKeyToFile(privateKeyFile, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(r.PrivateKey))
}

// Encrypt encrypts a message using RSA-OAEP.
func (r *RSAEncryptor) Encrypt(secretMessage string) string {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, r.PublicKey, []byte(secretMessage), label)
	utils.CheckError(err)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Decrypt decrypts a message using RSA-OAEP.
func (r *RSAEncryptor) Decrypt(cipherText string) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, r.PrivateKey, ct, label)
	utils.CheckError(err)
	fmt.Println("\nPlaintext:", string(plaintext))
	return string(plaintext)
}

// saveKeyToFile saves the key to a file.
func (r *RSAEncryptor) saveKeyToFile(filename, keyType string, keyData []byte) {
	file, err := os.Create(filename)
	utils.CheckError(err)
	defer file.Close()

	keyPEM := &pem.Block{Type: keyType, Bytes: keyData}
	err = pem.Encode(file, keyPEM)
	utils.CheckError(err)
}
