package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

var publicKey *rsa.PublicKey //nolint:gochecknoglobals // new for task

// const STEP = 400

func InitPublicKey(filePath string) error {
	publicKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	key, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}
	switch v := key.(type) {
	default:
		panic("unexpected key type")
	case *rsa.PublicKey:
		publicKey = v
	}
	return nil
}

func EncryptString(plaintext []byte) ([]byte, error) {
	msgLen := len(plaintext)
	// Не понял пока как подобрать число чтобы не было ошибки crypto/rsa: message too long for RSA key size.
	step := publicKey.Size() - 15 //nolint:gomnd // не понял пока как подобрать число
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		encryptedBlockBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext[start:finish])
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}
	return encryptedBytes, nil
}
