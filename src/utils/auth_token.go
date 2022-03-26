package utils

import (
	"encoding/json"
	"errors"
	fernet "github.com/fernet/fernet-go"
	"os"
	"time"
)

const (
	// The TTL of an AuthToken
	AuthTokenDuration = time.Hour * 24 * 30
)

// Represents an object for managing user authentication tokens.
type AuthToken struct {
	userId string
	email string
	secureText string
	expires time.Time
}

// Decodes an authentication token string into an AuthToken object.
func (*AuthToken) Decode(token string) (authToken *AuthToken, err error) {
	currentTime := time.Now()
	obj := make(map[string]string)
	obj["userId"] = ""
	obj["email"] = ""
	obj["secureText"] = ""
	obj["expires"] = ""
	key, err := fernet.DecodeKey(os.Getenv("APP_SECRET_KEY"))
	if err != nil {
		return nil, err
	}
	jsonTxt := fernet.VerifyAndDecrypt(
		[]byte(token),
		AuthTokenDuration,
		[]*fernet.Key{key},
	)
	err = json.Unmarshal(jsonTxt, &obj)
	if err != nil {
		return nil, err
	}
	isoTime, err := time.Parse(time.RFC3339, obj["expires"])
	if err != nil {
		return nil, err
	}
	if currentTime.After(isoTime.UTC()) {
		return nil, errors.New("token expired")
	}
	authToken = new(AuthToken)
	authToken.userId = obj["userId"]
	authToken.email = obj["email"]
	authToken.secureText = obj["secureText"]
	authToken.expires = isoTime.UTC()
	return authToken, nil
}

// Encodes an AuthToken object to an authentication token string.
func (authToken *AuthToken) Encode() (token string, err error) {
	currentTime := time.Now()
	currentTime.Add(AuthTokenDuration)
	obj := make(map[string]string)
	obj["userId"] = authToken.userId
	obj["email"] = authToken.email
	obj["secureText"] = authToken.secureText
	obj["expires"] = currentTime.UTC().Format(time.RFC3339)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	key, err := fernet.DecodeKey(os.Getenv("APP_SECRET_KEY"))
	if err != nil {
		return "", err
	}
	tok, err := fernet.EncryptAndSign(jsonBytes, key)
	if err != nil {
		return "", err
	}
	return string(tok), nil
}
