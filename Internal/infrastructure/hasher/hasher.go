package hasher

import (
	"crypto/sha256"
	errorsAPP "edtech/pkg/errors"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct{}

type HasherManager interface {
	Hash(password string) (string, error)
	CheckPassword(passhash string, password string) (bool, error)
	HashRefreshToken(token string) string
}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (b *Hasher) Hash(password string) (string, error) {

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("%w: %v", errorsAPP.ErrFailHashingPassword, err)
	}

	return string(passhash), nil
}

func (b *Hasher) CheckPassword(passhash string, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(passhash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, fmt.Errorf("%w: %w", errorsAPP.ErrFailCheckingPassword, err)
	}
	return true, nil
}

func (b *Hasher) HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
