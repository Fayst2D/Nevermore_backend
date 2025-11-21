package handler

import (
	"nevermore/internal/transport/handler/chat"

	"github.com/gin-gonic/gin"

	"nevermore/internal/service"
	"nevermore/internal/transport/handler/auth"
	"nevermore/internal/transport/handler/author"
	"nevermore/internal/transport/handler/book"
	"nevermore/internal/transport/handler/saved_author"
	"nevermore/internal/transport/handler/user"
	middleware2 "nevermore/internal/transport/middleware"
	tokenManager "nevermore/pkg/auth"
	"time"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	serv   service.Service
	router *gin.Engine
}

func New(serv service.Service, manager *tokenManager.Manager) *gin.Engine {
	handler := &Handler{
		serv:   serv,
		router: gin.Default(),
	}

	handler.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Разрешить все домены
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowHeaders:     []string{"*"}, // Разрешить все заголовки
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//добавление СВАГИ
	handler.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userHandler := user.New(serv)
	authorHandler := author.New(serv)
	savedAuthorHandler := saved_author.New(serv)
	bookHandler := book.New(serv)
	authHandler := auth.New(serv)
	chatHandler := chat.New(serv)

	handler.router.POST("/auth/register", authHandler.Register)
	handler.router.POST("/auth/login", authHandler.Login)

	protected := handler.router.Group("/")
	protected.Use(middleware2.AuthMiddleware(manager))
	protected.Use(middleware2.RateLimiter(1 * time.Second))
	{
		protected.GET("/user", userHandler.Get)
		protected.PUT("/user", userHandler.Update)
		protected.DELETE("/user", userHandler.Delete)

		protected.POST("/saved-author/create", savedAuthorHandler.Create)
		protected.DELETE("/saved-author/delete", savedAuthorHandler.Delete)
		protected.GET("/saved-author/list", savedAuthorHandler.GetList)

		protected.GET("/author/get/:id", authorHandler.Get)
		protected.GET("/author/list", authorHandler.GetAuthorsList)
		protected.POST("/author/update/:id", authorHandler.Update)
		protected.DELETE("/author/delete/:id", authorHandler.Delete)

		protected.POST("/book/upload", bookHandler.Create)

		protected.GET("/chat/ws", chatHandler.WebSocketHandler)
		protected.GET("/chat/messages", chatHandler.GetMessages)
		protected.GET("/chat/online", chatHandler.GetOnlineUsers)
		//личные сообщения
		protected.POST("/chat/private/send", chatHandler.SendPrivateMessage)
		protected.GET("/chat/private/conversation/:user_id", chatHandler.GetPrivateConversation)
		//protected.GET("/chat/private/conversations", chatHandler.GetConversations)
		protected.POST("/chat/private/mark-read", chatHandler.MarkMessagesAsRead)
		protected.GET("/chat/private/unread-count", chatHandler.GetUnreadCount)
	}

	return handler.router
}
