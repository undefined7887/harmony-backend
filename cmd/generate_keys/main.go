package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"os"

	"github.com/undefined7887/harmony-backend/internal/util/crypto"
)

func main() {
	if len(os.Args) < 4 {
		writeError(`usage:
	$ generate_keys <type (ed25519, ecdsa)> /path/to/output_private /path/to/output_public 
`)
	}

	switch os.Args[1] {
	case "ed25519":
		key := randomKeyEd25519()

		if err := cryptoutil.WritePrivateKey(os.Args[2], key); err != nil {
			writeError("failed to create ed25519 private key", err)
		}

		if err := cryptoutil.WritePublicKey(os.Args[3], key.Public()); err != nil {
			writeError("failed to create ed25519 public key", err)
		}

	case "ecdsa":
		key := randomKeyECDSA()

		if err := cryptoutil.WritePrivateKey(os.Args[2], key); err != nil {
			writeError("failed to create ECDSA private key", err)
		}

		if err := cryptoutil.WritePublicKey(os.Args[3], key.Public()); err != nil {
			writeError("failed to create ECDSA public key", err)
		}

	default:
		writeError("type %s not supported", os.Args[2])
	}

	log.Println("keys generated")
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
