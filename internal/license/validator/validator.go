package validator

import (
	"crypto"
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

// Validator is the struct that all funcs are attached to
type Validator struct {
	pk *rsa.PublicKey
}

// New returns a Validator or dies
func New() *Validator {
	r := strings.NewReader(getCert())
	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	dat, _ := pem.Decode(pemBytes)
	if dat == nil {
		log.Fatal("no PEM encoded key found")
	}

	if dat.Type != "PUBLIC KEY" {
		log.Fatalf("unknown  block type %s", dat.Type)
	}

	publicKey, err := x509.ParsePKIXPublicKey(dat.Bytes)
	if err != nil {
		log.Fatalf("public key error: %v", err)
	}

	k, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("not an RSA pubkey")
	}

	return &Validator{
		pk: k,
	}
}

// ValidateString validates a license string and checks the signature, returning either a license or an error.
func (v *Validator) ValidateString(arg string) (*format.Payload, error) {
	p := format.Payload{}
	if err := p.FromString(arg); err != nil {
		// log.WithFields(log.Fields{
		// 	"status": "failed",
		// }).Error("err")
		return nil, err
	}

	if err := v.ValidatePayload(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

// ValidatePayload accepts a payload and validates it, returning an error on failure
func (v *Validator) ValidatePayload(req *format.Payload) error {
	licenseBytes, err := json.Marshal(req.License)
	if err != nil {
		// log.WithFields(log.Fields{
		// 	"status": "failed",
		// }).Error("err")
		return err
	}

	// Create sha256 Hash
	hashed := sha256.Sum256(licenseBytes)

	// Compare the hash against the value in the signature
	if err := rsa.VerifyPKCS1v15(v.pk, crypto.SHA256, hashed[:], req.Signature); err != nil {
		// log.WithFields(log.Fields{
		// 	"status": "failed",
		// }).Error("err")
		return err
	}

	return nil
}
