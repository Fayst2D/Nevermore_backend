package handler

import (
	"nevermore/internal/service"
	"nevermore/internal/transport/handler/author"
	"nevermore/internal/transport/handler/saved_author"
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
	authorHandler := author.New(serv)
	savedAuthorHandler := saved_author.New(serv)

	protected := handler.router.Group("/")
	protected.Use(middleware2.RateLimiter(1 * time.Second))
	{
		protected.GET("/user/get", userHandler.Get)
		protected.POST("/user/update", userHandler.Update)
		protected.DELETE("/user/delete", userHandler.Delete)

		protected.POST("/saved-author/create", savedAuthorHandler.Create)
		protected.DELETE("/saved-author/delete", savedAuthorHandler.Delete)
		protected.GET("/saved-author/list", savedAuthorHandler.GetList)

		handler.router.GET("/author/get/:id", authorHandler.Get)
		handler.router.GET("/author/list", authorHandler.GetAuthorsList)
		handler.router.POST("/author/update/:id", authorHandler.Update)
		handler.router.DELETE("/author/delete/:id", authorHandler.Delete)
	}

	return handler.router
}
