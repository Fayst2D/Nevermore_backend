package handler

import (
	"nevermore/internal/service"
	"nevermore/internal/transport/handler/user"
	middleware2 "nevermore/internal/transport/middleware"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	serv   service.Service
	router *gin.Engine
}

func New(serv service.Service) *gin.Engine {
	handler := &Handler{
		serv:   serv,
		router: gin.New(),
	}

	//добавление СВАГИ
	handler.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userHandler := user.New(serv)

	protected := handler.router.Group("/")
	protected.Use(middleware2.RateLimiter(1 * time.Second))
	{
		protected.GET("/user/get", userHandler.Get)
		protected.POST("/user/update", userHandler.Update)
		protected.DELETE("/user/delete", userHandler.Delete)
	}

	return handler.router
}
