package issuer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"strings"

	"log"

	"github.com/DDLbots/escape-pod/internal/license/format"
)

// Issuer is the struct that all funcs are attached to
type Issuer struct {
	pk *rsa.PrivateKey
}

// New returns a licensor or dies
func New() *Issuer {
	r := strings.NewReader(getKey())
	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	dat, _ := pem.Decode(pemBytes)
	if dat == nil {
		log.Fatal("no PEM encoded key found")
	}

	if dat.Type != "PRIVATE KEY" && dat.Type != "RSA PRIVATE KEY" {
		log.Fatalf("unknown  block type %s", dat.Type)
	}

	pk, err := x509.ParsePKCS8PrivateKey(dat.Bytes)
	if err != nil {
		log.Fatalf("invalid license private key: %v", err)
	}

	return &Issuer{
		pk: pk.(*rsa.PrivateKey),
	}
}

// Generate returns a license key, or an error on failure.
func (l *Issuer) Generate(in *format.License) (string, error) {
	p := format.Payload{
		License: in,
	}

	licenseBytes, err := json.Marshal(p.License)
	if err != nil {
		return "", err
	}

	hashed := sha256.Sum256(licenseBytes)
	signature, err := rsa.SignPKCS1v15(rand.Reader, l.pk, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	p.Signature = signature

	return p.ToString()
}
