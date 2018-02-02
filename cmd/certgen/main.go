package main

import (
	"bytes"
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/jessevdk/go-flags"
)

// Options defines command line options
type Options struct {
	WorkDir string `short:"w" long:"work-dir" description:"Working directory"`
	LogFile string `short:"l" long:"log-file" description:"Log file output"`
	Verbose bool   `long:"verbose" description:"Enable verbose logging"`

	Validate       Validate
	GenerateCA     GenerateCA
	GenerateClient GenerateClient
	Revoke         Revoke
	Transcode      Transcode
	Daemon         Daemon
}

type certOpts struct {
	SerialNumber uint      `short:"s" long:"serial-number" description:"Client certificate serial number"`
	ValidBefore  time.Time `long:"valid-before" description:"Certificate validity start date"`
	ValidUntil   time.Time `long:"valid-until" description:"Certificate validity end date"`
	KeyUsage     []string  `long:"key-usage" description:"Certificate key uses"`
	ExtKeyUsage  []string  `long:"ext-key-usage" description:"Extended key uses"`

	IsCA    bool `long:"is-ca" description:"Is the certificate usable as a Certificate Authority"`
	PathLen uint `long:"path-len" description:"Certificate critical path length"`

	DNSNames       []string `long:"dns-names" description:"DNS names for certificate"`
	EmailAddresses []string `long:"email-addresses" description:"Email addresses for certificate"`
	IPAddresses    []net.IP `long:"ip-addresses" description:"IP addresses for certificate"`
}

// GenerateCA command generates a certificate authority
type GenerateCA struct {
	Name  flags.Filename `short:"n" long:"name" description:"Certificate authority name"`
	CAKey flags.Filename `short:"k" long:"ca-key" description:"Certificate authority private key file, if not provided a key will be generated"`
	certOpts
}

// DefaultGenerateCA creates the default GenerateCA configuration
func DefaultGenerateCA() GenerateClient {
	return GenerateClient{
		CAFile: "./ca.crt",
		CAKey:  "./ca.key",
		certOpts: certOpts{
			SerialNumber: 0,
			ValidBefore:  time.Now(),
			ValidUntil:   time.Now().Add(time.Hour * 24 * 365 * 10),
			KeyUsage:     []string{"DigitalSignature", "KeyEncipherment"},
			ExtKeyUsage:  []string{"ServerAuth", "ClientAuth"},
			IsCA:         true,
			PathLen:      1,
		},
	}
}

// GenerateClient command generates a client certificate and private key pair
type GenerateClient struct {
	Name      flags.Filename `short:"n" long:"name" description:"Client certificate name"`
	CAFile    flags.Filename `short:"c" long:"ca-file" description:"Certificate Authority certificate file"`
	CAKey     flags.Filename `short:"k" long:"ca-key" description:"Certificate Authority private key file"`
	ClientKey flags.Filename `short:"" long:"client-key" description:"Client certificate private key file, if not provided a key will be generated"`

	certOpts
}

// DefaultGenerateClient creates the default GenerateClient object
func DefaultGenerateClient() GenerateClient {
	return GenerateClient{
		CAFile: "./ca.crt",
		CAKey:  "./ca.key",
		certOpts: certOpts{
			SerialNumber: 0,
			ValidBefore:  time.Now(),
			ValidUntil:   time.Now().Add(time.Hour * 24 * 365 * 10),
			KeyUsage:     []string{"DigitalSignature", "KeyEncipherment"},
			ExtKeyUsage:  []string{"ServerAuth", "ClientAuth"},
			IsCA:         false,
			PathLen:      0,
		},
	}
}

// KeyUsages map to recognise usages for a CA cert
var KeyUsages = map[string]x509.KeyUsage{
	"DigitalSignature":  x509.KeyUsageDigitalSignature,
	"ContentCommitment": x509.KeyUsageContentCommitment,
	"KeyEncipherment":   x509.KeyUsageKeyEncipherment,
	"DataEncipherment":  x509.KeyUsageDataEncipherment,
	"KeyAgreement":      x509.KeyUsageKeyAgreement,
	"CertSign":          x509.KeyUsageCertSign,
	"CRLSign":           x509.KeyUsageCRLSign,
	"EncipherOnly":      x509.KeyUsageEncipherOnly,
	"DecipherOnly":      x509.KeyUsageDecipherOnly,
}

// ExtKeyUsages map for extended key usages
var ExtKeyUsages = map[string]x509.ExtKeyUsage{
	"Any":                        x509.ExtKeyUsageAny,
	"ServerAuth":                 x509.ExtKeyUsageServerAuth,
	"ClientAuth":                 x509.ExtKeyUsageClientAuth,
	"CodeSigning":                x509.ExtKeyUsageCodeSigning,
	"EmailProtection":            x509.ExtKeyUsageEmailProtection,
	"IPSECEndSystem":             x509.ExtKeyUsageIPSECEndSystem,
	"IPSECTunnel":                x509.ExtKeyUsageIPSECTunnel,
	"IPSECUser":                  x509.ExtKeyUsageIPSECUser,
	"TimeStamping":               x509.ExtKeyUsageTimeStamping,
	"OCSPSigning":                x509.ExtKeyUsageOCSPSigning,
	"MicrosoftServerGatedCrypto": x509.ExtKeyUsageMicrosoftServerGatedCrypto,
	"NetscapeServerGatedCrypto":  x509.ExtKeyUsageNetscapeServerGatedCrypto,
}

type Validate struct {
}

type Revoke struct {
}

type Transcode struct {
}

type Daemon struct {
}

func main() {
	fmt.Printf("CertGen")

	o := Options{}
	flags.Parse(&o)

}

func load(certFile, keyFile string) (*x509.Certificate, crypto.PrivateKey, error) {
	// Load certificate
	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, nil, err
	}

	// Load private key based on key type
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}

	key, err := x509.ParsePKIXPublicKey(keyData)

	return cert, key, err
}

func save(cert *x509.Certificate, key crypto.PrivateKey) error {

	certPem := bytes.NewBuffer(nil)
	err := pem.Encode(certPem, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if err != nil {
		return err
	}

	keyDer, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}

	keyType := ""
	switch key.(type) {
	case *rsa.PrivateKey:
		keyType = "RSA PRIVATE KEY"
	case *ecdsa.PrivateKey:
		keyType = "EC PRIVATE KEY"
	case *dsa.PrivateKey:
		keyType = "DSA PRIVATE KEY"
	}

	keyPem := bytes.NewBuffer(nil)
	err = pem.Encode(keyPem, &pem.Block{Type: keyType, Bytes: keyDer})
	if err != nil {
		return err
	}

	return nil
}
