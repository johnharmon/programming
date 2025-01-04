package main

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type JwtKey struct {
	KID             string `yaml:"key_id"`
	KeySecretString string `yaml:"key_secret"`
	KeySecret       []byte
}

func LoadJwtKeys(fileContent []byte, keyMap map[string]*JwtKey) (unmarshalErr error) {
	err := yaml.Unmarshal(fileContent, &keyMap)
	if err != nil {
		unmarshalErr = fmt.Errorf("error unmarshalling yaml into key map: %+v", err)
	}
	return unmarshalErr
}

func DecodeJwtSecrets(keyMap map[string]*JwtKey) {
	for key := range keyMap {
		if err := keyMap[key].GetSecret(); err != nil {
			fmt.Printf("Error decoding secret for key: %s\n", key)
		}
	}
}

func NewJwtKey(secret []byte) (key *JwtKey) {
	key.KID = uuid.NewString()
	key.KeySecret = secret
	key.KeySecretString = base64.StdEncoding.EncodeToString(secret)
	return key
}

func NewJwtKeyWithUUID(secret []byte, uuidString string) (key *JwtKey) {
	key.KID = uuidString
	key.KeySecret = secret
	key.KeySecretString = base64.StdEncoding.EncodeToString(secret)
	return key
}

// Assume yaml has been unmarshalled into a list of keys
func NewJwtKeyFromYaml(map[string]string) *JwtKey {
	return &JwtKey{}
}

func (k *JwtKey) GetSecret() (decodeErr error) {
	_, err := base64.StdEncoding.Decode(k.KeySecret, []byte(k.KeySecretString))
	if err != nil {
		decodeErr = fmt.Errorf("error decoding base64 secret: %+v", err)
	}
	return decodeErr
}
