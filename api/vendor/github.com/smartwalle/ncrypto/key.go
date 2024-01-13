package ncrypto

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type PrivateKey interface {
	Public() crypto.PublicKey
	Equal(crypto.PrivateKey) bool
}

type PublicKey interface {
	Equal(crypto.PublicKey) bool
}

type PrivateKeyDecoder []byte

func DecodePrivateKey(data []byte) PrivateKeyDecoder {
	return data
}

func (this PrivateKeyDecoder) decode() ([]byte, error) {
	if len(this) == 0 {
		return nil, errors.New("invalid private key")
	}

	if this[0] == '-' {
		block, _ := pem.Decode(this)
		if block == nil {
			return nil, errors.New("invalid private key")
		}
		return block.Bytes, nil
	}

	var data, err = base64decode(this)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this PrivateKeyDecoder) PKCS1() PKCS1PrivateKey {
	der, err := this.decode()
	if err != nil {
		return PKCS1PrivateKey{key: nil, err: err}
	}
	key, err := x509.ParsePKCS1PrivateKey(der)
	return PKCS1PrivateKey{key: key, err: err}
}

func (this PrivateKeyDecoder) PKCS8() PKCS8PrivateKey {
	der, err := this.decode()
	if err != nil {
		return PKCS8PrivateKey{key: nil, err: err}
	}
	key, err := x509.ParsePKCS8PrivateKey(der)
	return PKCS8PrivateKey{key: key, err: err}
}

type PKCS1PrivateKey struct {
	key *rsa.PrivateKey
	err error
}

func (this PKCS1PrivateKey) RSAPrivateKey() (*rsa.PrivateKey, error) {
	if this.err != nil {
		return nil, this.err
	}
	return this.key, nil
}

type PKCS8PrivateKey struct {
	key any
	err error
}

func (this PKCS8PrivateKey) PrivateKey() (any, error) {
	return this.key, this.err
}

func (this PKCS8PrivateKey) RSAPrivateKey() (*rsa.PrivateKey, error) {
	if this.err != nil {
		return nil, this.err
	}
	privateKey, ok := this.key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key is not a valid *rsa.PrivateKey")
	}
	return privateKey, nil
}

func (this PKCS8PrivateKey) ECDSAPrivateKey() (*ecdsa.PrivateKey, error) {
	if this.err != nil {
		return nil, this.err
	}
	privateKey, ok := this.key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("key is not a valid *ecdsa.PrivateKey")
	}
	return privateKey, nil
}

func (this PKCS8PrivateKey) ED25519PrivateKey() (*ed25519.PrivateKey, error) {
	if this.err != nil {
		return nil, this.err
	}
	privateKey, ok := this.key.(*ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("key is not a valid *ed25519.PrivateKey")
	}
	return privateKey, nil
}

type PublicKeyDecoder []byte

func DecodePublicKey(data []byte) PublicKeyDecoder {
	return data
}

func (this PublicKeyDecoder) decode() ([]byte, error) {
	if len(this) == 0 {
		return nil, errors.New("invalid private key")
	}

	if this[0] == '-' {
		block, _ := pem.Decode(this)
		if block == nil {
			return nil, errors.New("invalid public key")
		}
		return block.Bytes, nil
	}

	var data, err = base64decode(this)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this PublicKeyDecoder) PKCS1() PKCS1PublicKey {
	der, err := this.decode()
	if err != nil {
		return PKCS1PublicKey{key: nil, err: err}
	}
	key, err := x509.ParsePKCS1PublicKey(der)
	return PKCS1PublicKey{key: key, err: err}
}

func (this PublicKeyDecoder) PKIX() PKIXPublicKey {
	der, err := this.decode()
	if err != nil {
		return PKIXPublicKey{key: nil, err: err}
	}
	key, err := x509.ParsePKIXPublicKey(der)
	return PKIXPublicKey{key: key, err: err}
}

type PKCS1PublicKey struct {
	key *rsa.PublicKey
	err error
}

func (this PKCS1PublicKey) RSAPublicKey() (*rsa.PublicKey, error) {
	return this.key, this.err
}

type PKIXPublicKey struct {
	key any
	err error
}

func (this PKIXPublicKey) PublicKey() (any, error) {
	return this.key, this.err
}

func (this PKIXPublicKey) RSAPublicKey() (*rsa.PublicKey, error) {
	if this.err != nil {
		return nil, this.err
	}
	publicKey, ok := this.key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not a valid *rsa.PublicKey")
	}
	return publicKey, nil
}

func (this PKIXPublicKey) ECDSAPublicKey() (*ecdsa.PublicKey, error) {
	if this.err != nil {
		return nil, this.err
	}
	publicKey, ok := this.key.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not a valid *ecdsa.PublicKey")
	}
	return publicKey, nil
}

func (this PKIXPublicKey) ED25519PublicKey() (*ed25519.PublicKey, error) {
	if this.err != nil {
		return nil, this.err
	}
	publicKey, ok := this.key.(*ed25519.PublicKey)
	if !ok {
		return nil, errors.New("key is not a valid *ed25519.PublicKey")
	}
	return publicKey, nil
}

type PrivateKeyEncoder struct {
	key PrivateKey
}

func EncodePrivateKey(key PrivateKey) PrivateKeyEncoder {
	return PrivateKeyEncoder{key: key}
}

func (this PrivateKeyEncoder) PKCS1() ([]byte, error) {
	switch pri := this.key.(type) {
	case *rsa.PrivateKey:
		privateBytes := x509.MarshalPKCS1PrivateKey(pri)
		block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateBytes}

		var buffer bytes.Buffer
		if err := pem.Encode(&buffer, block); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", pri)
	}
}

func (this PrivateKeyEncoder) PKCS8() ([]byte, error) {
	privateBytes, err := x509.MarshalPKCS8PrivateKey(this.key)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{Type: "PRIVATE KEY", Bytes: privateBytes}

	var buffer bytes.Buffer
	if err = pem.Encode(&buffer, block); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

type PublicKeyEncoder struct {
	key PublicKey
}

func EncodePublicKey(key PublicKey) PublicKeyEncoder {
	return PublicKeyEncoder{key: key}
}

func (this PublicKeyEncoder) PKCS1() ([]byte, error) {
	switch pub := this.key.(type) {
	case *rsa.PublicKey:
		publicBytes := x509.MarshalPKCS1PublicKey(pub)
		block := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicBytes}

		var buffer bytes.Buffer
		if err := pem.Encode(&buffer, block); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported public key type: %T", pub)
	}
}

func (this PublicKeyEncoder) PKIX() ([]byte, error) {
	publicBytes, err := x509.MarshalPKIXPublicKey(this.key)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: publicBytes}

	var buffer bytes.Buffer
	if err = pem.Encode(&buffer, block); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func base64decode(data []byte) ([]byte, error) {
	var dBuf = make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(dBuf, data)
	return dBuf[:n], err
}
