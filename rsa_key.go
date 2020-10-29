package line_login

import (
	"github.com/5hields/line-login/linethrift"
	"crypto/rsa"
	"encoding/hex"
	"log"
	"math/big"
	"strconv"
)

func rsaKeyGen(key *linethrift.RSAKey) rsa.PublicKey {
	decN, err := hex.DecodeString(key.Nvalue)
	if err != nil {
		log.Fatal(err)
	}
	n := big.NewInt(0)
	n.SetBytes(decN)

	eVal, err := strconv.ParseInt(key.Evalue, 16, 32)
	if err != nil {
		log.Fatal(err)
	}
	e := int(eVal)
	returnKey := rsa.PublicKey{N: n, E: e}

	return returnKey
}
