package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role,omitempty"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID, role string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		Type:   accessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(getAccessTokenTTL())),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getJWTSecret()))
}

func GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	claims := TokenClaims{
		UserID: userID,
		Type:   refreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(getRefreshTokenTTL())),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID,
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getRefreshJWTSecret()))
}

func ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	return validateToken(tokenString, getJWTSecret(), accessTokenType)
}

func ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	return validateToken(tokenString, getRefreshJWTSecret(), refreshTokenType)
}

func GetRefreshTokenTTL() time.Duration {
	return getRefreshTokenTTL()
}

func validateToken(tokenString, secret, expectedType string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	if claims.Type != expectedType {
		return nil, errors.New("tipo de token inválido")
	}

	if claims.UserID == "" {
		return nil, errors.New("token sin user_id")
	}

	return claims, nil
}

func getJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

func getRefreshJWTSecret() string {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret != "" {
		return secret
	}
	return getJWTSecret()
}

func getAccessTokenTTL() time.Duration {
	minutes, err := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_TTL_MINUTES"))
	if err == nil && minutes > 0 {
		return time.Duration(minutes) * time.Minute
	}
	return 15 * time.Minute
}

func getRefreshTokenTTL() time.Duration {
	hours, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_TTL_HOURS"))
	if err == nil && hours > 0 {
		return time.Duration(hours) * time.Hour
	}
	return 7 * 24 * time.Hour
}

func generateJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
