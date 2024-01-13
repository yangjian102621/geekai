package ncrypto

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func DecodeCertificate(data []byte) (*x509.Certificate, error) {
	der, err := decodeCertificate(data)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func decodeCertificate(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("invalid certificate")
	}

	if src[0] == '-' {
		block, _ := pem.Decode(src)
		if block == nil {
			return nil, errors.New("invalid certificate")
		}
		return block.Bytes, nil
	}

	var data, err = base64decode(src)
	if err != nil {
		return nil, err
	}
	return data, nil
}
