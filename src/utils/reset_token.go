package utils

import (
	"encoding/json"
	"errors"
	fernet "github.com/fernet/fernet-go"
	"os"
	"time"
)

const (
	// The TTL of a ResetToken
	ResetTokenDuration = time.Hour * 5
)

// Represents an object for managing user account reset tokens.
type ResetToken struct {
	userId string
	email string
	message string
	expires time.Time
}

// Decodes a reset token string into a ResetToken object.
func (*ResetToken) Decode(token string) (resetToken *ResetToken, err error) {
	currentTime := time.Now()
	obj := make(map[string]string)
	obj["userId"] = ""
	obj["email"] = ""
	obj["message"] = ""
	obj["expires"] = ""
	key, err := fernet.DecodeKey(os.Getenv("APP_SECRET_KEY"))
	if err != nil {
		return nil, err
	}
	jsonTxt := fernet.VerifyAndDecrypt(
		[]byte(token),
		ResetTokenDuration,
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
	resetToken = new(ResetToken)
	resetToken.userId = obj["userId"]
	resetToken.email = obj["email"]
	resetToken.message = obj["message"]
	resetToken.expires = isoTime.UTC()
	return resetToken, nil
}

// Encodes a ResetToken object to a reset token string.
func (resetToken *ResetToken) Encode() (token string, err error) {
	currentTime := time.Now()
	currentTime.Add(ResetTokenDuration)
	obj := make(map[string]string)
	obj["userId"] = resetToken.userId
	obj["email"] = resetToken.email
	obj["message"] = resetToken.message
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
