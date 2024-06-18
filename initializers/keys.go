package initializers

import (
	"crypto/rand"
	"crypto/rsa"
)

var PrivateKey *rsa.PrivateKey
var PublicKey *rsa.PublicKey

func InitKeys() {
	var err error
	PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	PublicKey = &PrivateKey.PublicKey
}
