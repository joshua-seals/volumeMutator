package commands

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func GenerateTLSCerts(certPath string) (*bytes.Buffer, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2023),
		Subject: pkix.Name{
			Organization:  []string{"renci.org"},
			Country:       []string{"US"},
			Province:      []string{"North Carolina"},
			Locality:      []string{"Chapel Hill"},
			StreetAddress: []string{"Europa Center 100 Europa Drive, Suite 540"},
			PostalCode:    []string{"27517"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// CA private key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Println("Error: generating private key ", err)
		return nil, err
	}

	// Self signed CA certificate based on template above
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		log.Println("Error: generating self signed certificate ", err)
		return nil, err
	}

	// PEM encode CA certificate
	caPEM := new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	// caPrivKeyPEM := new(bytes.Buffer)
	// pem.Encode(caPrivKeyPEM, &pem.Block{
	// 	Type:  "RSA PRIVATE KEY",
	// 	Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	// })

	// Very important to mirror the name of the service
	// [name of service].[namespace].svc for dns
	// **** MAKE DYNAMIC
	dnsNames := []string{"volume-mutator-svc",
		"volume-mutator-svc.default", "volume-mutator-svc.default.svc"}

	// server cert config
	cert := &x509.Certificate{
		DNSNames:     dnsNames,
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName:    "webhookCert",
			Organization:  []string{"renci.org"},
			Country:       []string{"US"},
			Province:      []string{"North Carolina"},
			Locality:      []string{"Chapel Hill"},
			StreetAddress: []string{"Europa Center 100 Europa Drive, Suite 540"},
			PostalCode:    []string{"27517"},
		},
		// Ensure valid at localhost too
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// server private key
	serverPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Println("Error: generating server priv key ", err)
		return nil, err
	}
	// sign the server certificate, note parent is ca created at the beginning
	serverCertBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &serverPrivKey.PublicKey, serverPrivKey)
	//serverCertBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &serverPrivKey.PublicKey, caPrivKey)
	if err != nil {
		log.Println("Error: creating server cert ", err)
		return nil, err
	}
	// PEM encode the server cert and key
	serverCertPEM := new(bytes.Buffer)
	_ = pem.Encode(serverCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertBytes,
	})

	serverPrivKeyPEM := new(bytes.Buffer)
	_ = pem.Encode(serverPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverPrivKey),
	})

	err = os.MkdirAll(certPath, 0755)
	if err != nil {
		log.Println("Error: Creating Directory ", err)
		return nil, err
	}

	err = WriteFile(certPath+"tls.crt", serverCertPEM)
	if err != nil {
		log.Println("Error: Writing tls.crt ", err)
		return nil, err
	}
	err = WriteFile(certPath+"tls.key", serverPrivKeyPEM)
	if err != nil {
		log.Println("Error: Writing tls.key ", err)
		return nil, err

	}
	return caPEM, nil
}

// WriteFile writes data in the file at the given path
func WriteFile(filepath string, sCert *bytes.Buffer) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(sCert.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Ref: https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
