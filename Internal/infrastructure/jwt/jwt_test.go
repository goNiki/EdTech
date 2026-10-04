package jwt_test

import (
	"testing"
	"time"

	"edtech/internal/domain"
	jwtinfra "edtech/internal/infrastructure/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockJWTConfig struct {
	secret     string
	accessExp  time.Duration
	refreshExp time.Duration
}

func (m *mockJWTConfig) Secret() string              { return m.secret }
func (m *mockJWTConfig) AccessExp() time.Duration    { return m.accessExp }
func (m *mockJWTConfig) RefreshExp() time.Duration   { return m.refreshExp }

func TestJwtManager_ParseToken_BestPractices(t *testing.T) {
	secret := "test-secret-key-32-bytes-long!!"
	cfg := &mockJWTConfig{
		secret:     secret,
		accessExp:  15 * time.Minute,
		refreshExp: 24 * time.Hour,
	}
	manager := jwtinfra.NewJwtManager(cfg)

	user := &domain.User{
		ID:       42,
		Username: "alice",
		Role:     domain.RoleStudent,
	}

	t.Run("valid access token parses successfully", func(t *testing.T) {
		tokenStr, err := manager.GenerateAccessToken(user, time.Now())
		require.NoError(t, err)

		id, role, err := manager.ParseToken(tokenStr)
		require.NoError(t, err)
		assert.Equal(t, int64(42), id)
		assert.Equal(t, string(domain.RoleStudent), role)
	})

	t.Run("token without exp claim is rejected", func(t *testing.T) {
		claims := jwtinfra.ClaimsAccessToken{
			ID:       user.ID,
			Username: user.Username,
			Role:     string(user.Role),
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt: jwt.NewNumericDate(time.Now()),
				Issuer:   "auth",
				// ExpiresAt omitted intentionally
			},
		}
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := tok.SignedString([]byte(secret))
		require.NoError(t, err)

		_, _, err = manager.ParseToken(tokenStr)
		require.Error(t, err, "must reject token without exp claim")
	})

	t.Run("token with invalid issuer is rejected", func(t *testing.T) {
		claims := jwtinfra.ClaimsAccessToken{
			ID:       user.ID,
			Username: user.Username,
			Role:     string(user.Role),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "evil-issuer",
			},
		}
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := tok.SignedString([]byte(secret))
		require.NoError(t, err)

		_, _, err = manager.ParseToken(tokenStr)
		require.Error(t, err, "must reject token with untrusted issuer")
	})

	t.Run("expired token is rejected", func(t *testing.T) {
		expiredTime := time.Now().Add(-1 * time.Hour)
		claims := jwtinfra.ClaimsAccessToken{
			ID:       user.ID,
			Username: user.Username,
			Role:     string(user.Role),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiredTime.Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(expiredTime),
				Issuer:    "auth",
			},
		}
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := tok.SignedString([]byte(secret))
		require.NoError(t, err)

		_, _, err = manager.ParseToken(tokenStr)
		require.Error(t, err, "must reject expired token")
	})
}
