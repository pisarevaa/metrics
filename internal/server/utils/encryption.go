package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

var privateKey *rsa.PrivateKey //nolint:gochecknoglobals // new for task

func InitPrivateKey(filePath string) error {
	privateKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}
	return nil
}

func DecryptString(ciphertext []byte) ([]byte, error) {
	msgLen := len(ciphertext)
	var decryptedBytes []byte
	step := privateKey.PublicKey.Size()

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBlockBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext[start:finish])
		if err != nil {
			// panic(err)
			return nil, err
		}
		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}
	return decryptedBytes, nil
}
