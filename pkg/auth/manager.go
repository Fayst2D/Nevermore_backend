package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
	GenerateTokenPair(userID string) (*Token, error)
	RefreshAccessToken(refreshTokenString string) (*Token, error)
}

type Manager struct {
	signingKey string
}

type CustomClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

// RefreshTokenClaims - кастомные claims для refresh токена
type RefreshTokenClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) GenerateTokenPair(userId string) (*Token, error) {
	// Access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			Issuer:    "your-app",
		},
	})

	accessTokenString, err := accessToken.SignedString([]byte(m.signingKey))
	if err != nil {
		return nil, err
	}

	// Refresh token с добавлением userID
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshTokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Issuer:    "your-app",
		},
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(m.signingKey))
	if err != nil {
		return nil, err
	}

	token := &Token{
		At: accessTokenString,
		Rt: refreshTokenString,
	}

	return token, nil
}

// Обновление access token с помощью refresh token
func (m *Manager) RefreshAccessToken(refreshTokenString string) (*Token, error) {
	// Проверяем refresh token с кастомными claims
	token, err := jwt.ParseWithClaims(refreshTokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Извлекаем claims и получаем userID из refresh токена
	if claims, ok := token.Claims.(*RefreshTokenClaims); ok {
		// Генерируем новую пару токенов с извлеченным userID
		return m.GenerateTokenPair(claims.UserId)
	}

	return nil, fmt.Errorf("failed to extract claims from refresh token")
}

// Реализация остальных методов интерфейса
func (m *Manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserId, nil
	}

	return "", fmt.Errorf("invalid access token")
}

func (m *Manager) NewRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(720 * time.Hour).Unix(), // 30 дней
			IssuedAt:  time.Now().Unix(),
		},
	})

	return token.SignedString([]byte(m.signingKey))
}
