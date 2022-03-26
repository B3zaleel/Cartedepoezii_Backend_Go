package utils
// Sourced from: https://golangcode.com/argon2-password-hashing/

import (
    "os"
    "fmt"
    "strings"
	"crypto/rand"
    "crypto/subtle"
    "encoding/base64"

	"golang.org/x/crypto/argon2"
)

// Represents configurations for hashing.
type HashConfig struct {
    time uint32
	memory uint32
	threads uint8
	keyLen uint32
}

// Generates a hash from the given password.
func GenerateHash(password string) (string, error) {
    hashCfg := &HashConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}
	salt := []byte(os.Getenv("PWD_SALT"))
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey(
        []byte(password),
        salt,
        hashCfg.time,
        hashCfg.memory,
        hashCfg.threads,
        hashCfg.keyLen,
    )
	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	pwdHash := fmt.Sprintf(
        "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
        argon2.Version,
        hashCfg.memory,
        hashCfg.time,
        hashCfg.threads,
        b64Salt,
        b64Hash,
    )
	return pwdHash, nil
}

// Checks if the given password matches the given hash.
func IsValidPassword(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	hashCfg := &HashConfig{}
	_, err := fmt.Sscanf(
        parts[3],
        "m=%d,t=%d,p=%d",
        &hashCfg.memory,
        &hashCfg.time,
        &hashCfg.threads,
    )
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	hashCfg.keyLen = uint32(len(decodedHash))
	comparisonHash := argon2.IDKey(
        []byte(password),
        salt,
        hashCfg.time,
        hashCfg.memory,
        hashCfg.threads,
        hashCfg.keyLen,
    )
	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}
