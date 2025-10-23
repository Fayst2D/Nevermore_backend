package middleware

import (
	"net/http"
	"nevermore/pkg/logger"
	"strings"

	"nevermore/pkg/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenManager *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		log := logger.Get()
		log.Info().Msgf("auth header: %v", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		log.Info().Msgf("Token: %s", tokenString)

		// Парсим и проверяем токен
		userID, err := tokenManager.Parse(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Сохраняем userID в контекст Gin для использования в обработчиках
		c.Set("userID", userID)

		// Продолжаем выполнение
		c.Next()
	}
}
