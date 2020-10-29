package line_login

import (
	"github.com/5hields/line-login/linethrift"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

func setPassword(msg string) string {
	sep := rune(len(msg))
	return string(sep)
}

func encryptCredential(email, pwd string, key *linethrift.RSAKey) string {
	rsaKey := rsaKeyGen(key)
	msg := []byte(setPassword(key.SessionKey) + key.SessionKey + setPassword(email) + email + setPassword(pwd) + pwd)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, &rsaKey, msg)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(cipherText)
}

func xor(data []byte) (r []byte) {
	length := len(data) / 2
	r = make([]byte, length)
	for i := 0; i < length; i++ {
		r[i] = data[i] ^ data[length+i]
	}
	return
}

func getSHA256Sum(data ...[]byte) (r []byte) {
	sha := sha256.New()
	for _, update := range data {
		sha.Write(update)
	}
	r = sha.Sum(nil)
	return
}

func aesECBEncrypt(aesKey, plainText []byte) (cipherText []byte, err error) {
	cipherText = make([]byte, len(plainText))
	block, err := aes.NewCipher(aesKey)
	ecb := newECBEncrypter(block)
	ecb.CryptBlocks(cipherText, plainText)
	return
}

func getHashKeyChain(sharedSecret, encKeyChain []byte) []byte {
	hashKey := getSHA256Sum(sharedSecret, []byte("Key"))
	data := xor(getSHA256Sum(encKeyChain))
	encKey, _ := aesECBEncrypt(hashKey, data)
	return encKey
}

func createSecret(data []byte) (string, []byte) {
	pin, _ := createPinCode()
	enc, _ := aesECBEncrypt(getSHA256Sum([]byte(pin)), data)
	return pin, enc
}

func createPinCode() (string, error) {
	pin, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return "", err
	}
	// padding
	return fmt.Sprintf("%06d", pin), nil
}
