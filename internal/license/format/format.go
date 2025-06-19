package format

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// Payload is the format for the license + wrapper
type Payload struct {
	License   *License `json:"payload,omitempty" bson:"license"`
	Signature []byte   `json:"signature,omitempty" bson:"signature"`
}

// ToString returns the key string
func (l *Payload) ToString() (string, error) {
	payloadBytes, err := json.Marshal(l)
	if err != nil {
		return "", err
	}

	key := base64.StdEncoding.EncodeToString(payloadBytes)

	return key, nil
}

// FromString populates a payload from a string, or errors on failure
func (l *Payload) FromString(in string) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return errors.New("invalid license key")
	}

	if err := json.Unmarshal(b, l); err != nil {
		return errors.New("invalid license key")
	}
	return nil
}

// License is the license format
type License struct {
	Email   string `json:"email,omitempty" bson:"email"`
	Version string `json:"version,omitempty" bson:"version"`
	Bot     string `json:"bot,omitempty" bson:"bot"`
}
