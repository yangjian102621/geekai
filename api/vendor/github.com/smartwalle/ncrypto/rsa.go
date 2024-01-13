package ncrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

func split(data []byte, max int) [][]byte {
	var src = make([]byte, len(data))
	copy(src, data)

	var chunks = make([][]byte, 0, (len(src)%max)+1)
	if len(src) <= max {
		return append(chunks, src)
	}
	for len(src) >= max {
		chunks = append(chunks, src[:max])
		src = src[max:]
	}
	if len(src) > 0 {
		chunks = append(chunks, src)
	}
	return chunks
}

// RSAEncrypt 使用公钥 key 对数据 plaintext 进行加密
func RSAEncrypt(plaintext []byte, key *rsa.PublicKey) ([]byte, error) {
	var chunks = split(plaintext, key.N.BitLen()/8-11)
	var ciphertext = make([]byte, 0, 0)

	for _, chunk := range chunks {
		var data, err = rsa.EncryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return nil, err
		}
		ciphertext = append(ciphertext, data...)
	}

	return ciphertext, nil
}

// RSADecrypt 使用私钥 key 对数据 ciphertext 进行解密
func RSADecrypt(ciphertext []byte, key *rsa.PrivateKey) ([]byte, error) {
	var chunks = split(ciphertext, key.PublicKey.N.BitLen()/8)
	var plaintext = make([]byte, 0, 0)

	for _, chunk := range chunks {
		var data, err = rsa.DecryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return nil, err
		}
		plaintext = append(plaintext, data...)
	}
	return plaintext, nil
}

func RSASignPKCS1v15(plaintext []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(plaintext)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func RSAVerifyPKCS1v15(ciphertext, signature []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(ciphertext)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, signature)
}

func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return private, &private.PublicKey, err
}
