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
	UserId string
	Email string
	Message string
	Expires time.Time
}

// Decodes a reset token string into a ResetToken object.
func DecodeResetToken(token string) (resetToken *ResetToken, err error) {
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
	resetToken.UserId = obj["userId"]
	resetToken.Email = obj["email"]
	resetToken.Message = obj["message"]
	resetToken.Expires = isoTime.UTC()
	return resetToken, nil
}

// Encodes a ResetToken object to a reset token string.
func EncodeResetToken(resetToken *ResetToken) (token string, err error) {
	currentTime := time.Now()
	currentTime.Add(ResetTokenDuration)
	obj := make(map[string]string)
	obj["userId"] = resetToken.UserId
	obj["email"] = resetToken.Email
	obj["message"] = resetToken.Message
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
