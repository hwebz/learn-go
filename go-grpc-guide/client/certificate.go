package client

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

const (
	CertificatePassphrase = "hadm@123"
)

func LoadX509KeyPair(certFile, keyFile, passphrase string) (tls.Certificate, error) {
	certPEMBlock, err := ioutil.ReadFile(certFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to read certificate file: %w", err)
	}

	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to read key file: %w", err)
	}

	cert := tls.Certificate{}
	cert.Certificate = append(cert.Certificate, certPEMBlock)

	var keyDERBlock *pem.Block
	for {
		keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock)
		if keyDERBlock == nil {
			return tls.Certificate{}, fmt.Errorf("failed to decode PEM block containing the key")
		}
		if x509.IsEncryptedPEMBlock(keyDERBlock) {
			decryptedKeyDERBlock, err := x509.DecryptPEMBlock(keyDERBlock, []byte(passphrase))
			if err != nil {
				return tls.Certificate{}, fmt.Errorf("failed to decrypt PEM block: %w", err)
			}
			keyDERBlock.Bytes = decryptedKeyDERBlock
			keyDERBlock.Headers = nil
		} else {
			// Handle PKCS#8 encrypted key if DEK-Info is not present
			if keyDERBlock.Type == "ENCRYPTED PRIVATE KEY" {
				privateKey, err := x509.DecryptPEMBlock(keyDERBlock, []byte(passphrase))
				if err != nil {
					return tls.Certificate{}, fmt.Errorf("failed to decrypt PKCS#8 private key: %w", err)
				}
				keyDERBlock.Bytes = privateKey
			}
		}

		cert.PrivateKey, err = parsePrivateKey(keyDERBlock.Bytes)
		if err != nil {
			return tls.Certificate{}, fmt.Errorf("failed to parse private key: %w", err)
		}
		break
	}

	return cert, nil
}

func parsePrivateKey(der []byte) (interface{}, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("unknown private key type or invalid format")
}
