package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/undefined7887/harmony-backend/internal/util/crypto"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		writeError(`usage:
	$ generate_keys /path/to/output <type (ed25519, ecdsa)>
`)
	}

	switch os.Args[2] {
	case "ed25519":
		if err := cryptoutil.WritePrivateKey(os.Args[1], randomKeyEd25519()); err != nil {
			writeError("failed to create ed25519 private key", err)
		}

	case "ecdsa":
		if err := cryptoutil.WritePrivateKey(os.Args[1], randomKeyECDSA()); err != nil {
			writeError("failed to create ECDSA private key", err)
		}

	default:
		writeError("type %s not supported", os.Args[2])
	}

	log.Println("private key generated")
}

func randomKeyEd25519() ed25519.PrivateKey {
	_, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		writeError("failed to generate ed25519 keys %v", err)
	}

	return private
}

func randomKeyECDSA() *ecdsa.PrivateKey {
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		writeError("failed to generate ECDSA keys %v", err)
	}

	return private
}

func writeError(format string, v ...any) {
	log.Printf(format+"\n", v...)
	os.Exit(1)
}
