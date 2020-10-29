package line_login

import (
	"crypto/rand"
	"golang.org/x/crypto/curve25519"
	"log"
)

type e2eeKeyPair struct {
	Pri *[32]byte
	Pub *[32]byte
}

func genKey() (privateKey, publicKey *[32]byte) {
	privateKey, publicKey = new([32]byte), new([32]byte)
	if _, err := rand.Read(privateKey[:]); err != nil {
		log.Fatal(err)
	}

	curve25519.ScalarBaseMult(publicKey, privateKey)
	return
}

func newE2EEKeyPair() (r *e2eeKeyPair) {
	priKey, pubKey := genKey()
	r = &e2eeKeyPair{
		Pri: priKey, Pub: pubKey,
	}
	return
}

func (e2ee e2eeKeyPair) GenSharedSecret(keyDate []byte) []byte {
	sharedSecret, err := curve25519.X25519(e2ee.Pri[:], keyDate)
	if err != nil {
		log.Fatal(err)
	}
	return sharedSecret
}
