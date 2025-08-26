package hasher

import (
	errorsAPP "edtech/pkg/errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (b *BcryptHasher) Hash(password string) (string, error) {

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("%w: %v", errorsAPP.ErrFailHashingPassword, err)
	}

	return string(passhash), nil
}
