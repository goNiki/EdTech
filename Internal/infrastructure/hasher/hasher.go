package hasher

import (
	"crypto/sha256"
	errorsAPP "edtech/pkg/errors"
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct{}

type HasherManager interface {
	Hash(password string) (string, error)
	CheckPassword(passhash string, password string) bool
	HashRefreshToken(token string) (string, error)
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

func (b *Hasher) CheckPassword(passhash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(passhash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false
		}
		log.Fatal("error check Password Internal")
		return false
	}
	return true
}

func (b *Hasher) HashRefreshToken(token string) (string, error) {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:]), nil
}
