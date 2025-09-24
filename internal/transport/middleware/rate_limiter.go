package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter создает middleware для ограничения частоты запросов
func RateLimiter(duration time.Duration) gin.HandlerFunc {
	limits := make(map[string]time.Time)
	var mu sync.Mutex

	return func(c *gin.Context) {
		// Получаем userID из контекста (установленного в AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Преобразуем userID в строку (предполагая, что это string)
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(500, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		mu.Lock()
		lastTime, exists := limits[userIDStr]
		now := time.Now()

		if exists && now.Sub(lastTime) < duration {
			mu.Unlock()
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		limits[userIDStr] = now
		mu.Unlock()

		c.Next()
	}
}
