package users

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
)

// ValidateAndHashPwd validates and hashes a password.
func ValidateAndHashPwd(password string, minimumLength uint) (string, error) {
	if uint(len(password)) < minimumLength {
		return "", fberrors.ErrShortPassword{MinimumLength: minimumLength}
	}

	if _, ok := commonPasswords[password]; ok {
		return "", fberrors.ErrEasyPassword
	}

	return HashPwd(password)
}

// ValidateMinLengthAndHashPwd checks only the minimum length (skipping the
// common-password list) and hashes the password. Admins use it to set a weak
// password on purpose.
func ValidateMinLengthAndHashPwd(password string, minimumLength uint) (string, error) {
	if uint(len(password)) < minimumLength {
		return "", fberrors.ErrShortPassword{MinimumLength: minimumLength}
	}
	return HashPwd(password)
}

// HashPwd hashes a password.
func HashPwd(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPwd checks if a password is correct.
func CheckPwd(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RandomPwd returns a random password made from passwordLength random bytes,
// base64 (URL alphabet) encoded so it can be typed. It is used for accounts
// that never log in with a password, such as proxy-auth users.
func RandomPwd(passwordLength uint) (string, error) {
	randomPasswordBytes := make([]byte, passwordLength)
	var _, err = rand.Read(randomPasswordBytes)
	if err != nil {
		return "", err
	}

	// This is done purely to make the password human-readable
	var randomPasswordString = base64.URLEncoding.EncodeToString(randomPasswordBytes)
	return randomPasswordString, nil
}
