package jwt

import (
	"edtech/internal/domain"
	"edtech/internal/infrastructure/config"
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	Secret     string
	AccessExp  time.Duration
	RefreshExp time.Duration
}

type TokenManager interface {
	GenerateAccessToken(u *domain.User) (string, error)
	GenerateRefreshToken(id int64) (string, time.Time, error)
	GetAccessTokenExpiresIn() int64
	ParseToken(tokenStr string) (int64, string, error)
}

type ClaimsAccessToken struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
type ClaimsRefreshToken struct {
	ID int64 `json:"id"`
	jwt.RegisteredClaims
}

func NewJwtManager(cfg *config.JWTConfig) *JwtManager {
	return &JwtManager{
		Secret:     cfg.Secret,
		AccessExp:  cfg.AccesExp,
		RefreshExp: cfg.RefreshExp,
	}
}

func (j *JwtManager) GenerateAccessToken(u *domain.User) (string, error) {
	claims := ClaimsAccessToken{
		ID:       u.ID,
		Username: u.Username,
		Role:     string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstring, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", fmt.Errorf("%w, %v", errorsAPP.ErrFailedCreateJWT, err)
	}

	return tokenstring, nil
}

func (j *JwtManager) GenerateRefreshToken(id int64) (string, time.Time, error) {
	expiresAt := time.Now().Add(j.RefreshExp)
	claims := ClaimsRefreshToken{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstring, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", time.Now(), fmt.Errorf("%w, %v", errorsAPP.ErrFailedCreateRefreshJwt, err)
	}

	return tokenstring, expiresAt, nil
}

func (j *JwtManager) GetAccessTokenExpiresIn() int64 {
	return int64(j.AccessExp.Seconds())
}

func (j *JwtManager) ParseToken(tokenStr string) (int64, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ClaimsAccessToken{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("%w, %v", errorsAPP.ErrInvalidJWT, err)
	}

	claims, ok := token.Claims.(*ClaimsAccessToken)

	if !ok || !token.Valid {
		return 0, "", errorsAPP.ErrInvalidJWT
	}

	return claims.ID, claims.Role, nil

}
