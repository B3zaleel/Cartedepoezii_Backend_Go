package utils

import (
	"encoding/json"
	"errors"
	"fernet"
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
	jsonTxt := fernet.VerifyAndDecrypt(token, ResetTokenDuration, key)
	err = json.Unmarshal(jsonTxt, obj)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(LayoutISO, obj["expires"])
	if err != nil {
		return nil, err
	}
	if currentTime.After(t.UTC()) {
		return nil, errors.New("Token expired.")
	}
	resetToken = new(ResetToken)
	resetToken.userId = obj["userId"]
	resetToken.email = obj["email"]
	resetToken.message = obj["message"]
	resetToken.expires = t.UTC()
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
	obj["expires"] = currentTime.UTC().Format(LayoutISO)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	key, err := fernet.DecodeKey(os.Getenv("APP_SECRET_KEY"))
	if err != nil {
		return "", err
	}
	token, err = fernet.EncryptAndSign(jsonBytes, key)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
