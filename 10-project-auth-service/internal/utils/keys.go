package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
)

func EnsureRSAKeys(privatePath, publicPath string) error {
	// Jika file private key dan public key sudah ada, lewati generator
	if _, err := os.Stat(privatePath); err == nil {
		if _, err := os.Stat(publicPath); err == nil {
			return nil
		}
	}

	// Buat parent directory jika belum ada
	if err := os.MkdirAll(filepath.Dir(privatePath), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(publicPath), 0755); err != nil {
		return err
	}

	// Generate RSA 2048-bit key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 1. Simpan Private Key ke format PEM
	privateFile, err := os.Create(privatePath)
	if err != nil {
		return err
	}
	defer privateFile.Close()

	privatePEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateFile, privatePEM); err != nil {
		return err
	}

	// 2. Simpan Public Key ke format PEM
	publicFile, err := os.Create(publicPath)
	if err != nil {
		return err
	}
	defer publicFile.Close()

	publicDer, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	publicPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicDer,
	}
	if err := pem.Encode(publicFile, publicPEM); err != nil {
		return err
	}

	return nil
}
