package nsign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

type RSAMethod struct {
	h          crypto.Hash
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRSAMethod(h crypto.Hash, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *RSAMethod {
	var nRSA = &RSAMethod{}
	nRSA.h = h
	nRSA.privateKey = privateKey
	nRSA.publicKey = publicKey
	return nRSA
}

func (this *RSAMethod) Sign(data []byte) ([]byte, error) {
	var h = this.h.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, this.privateKey, this.h, hashed)
}

func (this *RSAMethod) Verify(data []byte, signature []byte) error {
	var h = this.h.New()
	if _, err := h.Write(data); err != nil {
		return err
	}
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(this.publicKey, this.h, hashed, signature)
}
