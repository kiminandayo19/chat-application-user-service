package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTHelper struct {
	secretKey []byte
	exp       time.Duration
}

func NewJWTHelper(secretKey string, exp time.Duration) *JWTHelper {
	return &JWTHelper{
		secretKey: []byte(secretKey),
		exp:       exp,
	}
}

func (h *JWTHelper) GenerateToken(userId, username string) (string, error) {
	claims := &JWTClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(h.secretKey)

	if err != nil {
		return "", err
	}

	return signed, nil
}

func (h *JWTHelper) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return h.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to parse token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token claims")
	}

	return claims, nil
}

func (h *JWTHelper) RefreshToken(oldToken string) (string, error) {
	claims, err := h.ValidateToken(oldToken)
	if err != nil {
		return "", err
	}
	return h.GenerateToken(claims.UserId, claims.Username)
}

func ExtractBearerToken(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("Invalid authorization format")
	}

	return authHeader[:7], nil
}
