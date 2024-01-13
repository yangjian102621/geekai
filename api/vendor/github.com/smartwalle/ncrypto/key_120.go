//go:build go1.20.0

package ncrypto

import (
	"crypto/ecdh"
	"errors"
)

func (this PKIXPublicKey) ECDHPublicKey() (*ecdh.PublicKey, error) {
	if this.err != nil {
		return nil, this.err
	}
	publicKey, ok := this.key.(*ecdh.PublicKey)
	if !ok {
		return nil, errors.New("key is not a valid *ecdh.PublicKey")
	}
	return publicKey, nil
}

func (this PKCS8PrivateKey) ECDHPrivateKey() (*ecdh.PrivateKey, error) {
	if this.err != nil {
		return nil, this.err
	}
	privateKey, ok := this.key.(*ecdh.PrivateKey)
	if !ok {
		return nil, errors.New("key is not a valid *ecdh.PrivateKey")
	}
	return privateKey, nil
}
