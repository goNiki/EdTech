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
	GenerateAccessToken(u *domain.User, now time.Time) (string, error)
	GenerateRefreshToken(id int64, now time.Time) (string, time.Time, error)
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

func NewJwtManager(cfg config.JWT) *JwtManager {
	return &JwtManager{
		Secret:     cfg.Secret(),
		AccessExp:  cfg.AccessExp(),
		RefreshExp: cfg.RefreshExp(),
	}
}

func (j *JwtManager) GenerateAccessToken(u *domain.User, now time.Time) (string, error) {
	claims := ClaimsAccessToken{
		ID:       u.ID,
		Username: u.Username,
		Role:     string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.AccessExp)),
			IssuedAt:  jwt.NewNumericDate(now),
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

func (j *JwtManager) GenerateRefreshToken(id int64, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(j.RefreshExp)
	claims := ClaimsRefreshToken{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstring, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", now, fmt.Errorf("%w, %v", errorsAPP.ErrFailedCreateRefreshJwt, err)
	}

	return tokenstring, expiresAt, nil
}

func (j *JwtManager) GetAccessTokenExpiresIn() int64 {
	return int64(j.AccessExp.Seconds())
}

func (j *JwtManager) ParseToken(tokenStr string) (int64, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ClaimsAccessToken{}, func(t *jwt.Token) (any, error) {
		return []byte(j.Secret), nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer("auth"),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return 0, "", fmt.Errorf("%w, %v", errorsAPP.ErrInvalidJWT, err)
	}

	claims, ok := token.Claims.(*ClaimsAccessToken)

	if !ok || !token.Valid {
		return 0, "", errorsAPP.ErrInvalidJWT
	}

	return claims.ID, claims.Role, nil
}
