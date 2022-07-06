package app

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

var (
	CertFile string = "localhost.crt"
	KeyFile  string = "localhost.key"
)

func GenerateCert() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		return err
	}

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	err = os.WriteFile(KeyFile, privateKeyPEM.Bytes(), 0600)
	if err != nil {
		return err
	}

	myCert := &x509.Certificate{
		SerialNumber: big.NewInt(2022),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"UZ"},
			Province:      []string{""},
			Locality:      []string{"Tashkent"},
			StreetAddress: []string{"Navoi street"},
			PostalCode:    []string{"10015"},
		},
		IPAddresses: []net.IP{{127, 0, 0, 1}, net.IPv6loopback},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(5, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
		DNSNames:    []string{"localhost"},
	}

	pemContent, err := ioutil.ReadFile(KeyFile)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(pemContent)

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	myCertBytes, err := x509.CreateCertificate(rand.Reader, myCert, myCert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: myCertBytes,
	})

	err = os.WriteFile(CertFile, certPEM.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
