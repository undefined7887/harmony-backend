package cryptoutil

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

const (
	TokenSize = 64
)

func Token() string {
	result := make([]byte, TokenSize)

	if _, err := rand.Read(result); err != nil {
		panic(fmt.Sprintf("unexpected error during token generation: %v", err))
	}

	return hex.EncodeToString(result)
}

func ReadPrivateKey(path string) (crypto.Signer, error) {
	privateKeyPem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	privateKeyDer, _ := pem.Decode(privateKeyPem)
	if privateKeyDer == nil {
		return nil, errors.New("cannot decode pem")
	}

	if privateKeyDer.Type != "PRIVATE KEY" {
		return nil, errors.New("key is not a private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyDer.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey.(crypto.Signer), nil
}

func WritePublicKey(path string, key any) error {
	publicKeyDer, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}

	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDer,
	})

	return os.WriteFile(path, publicKeyPem, os.ModePerm)
}

func WritePrivateKey(path string, key any) error {
	privateKeyDer, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}

	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDer,
	})

	return os.WriteFile(path, privateKeyPem, os.ModePerm)
}
