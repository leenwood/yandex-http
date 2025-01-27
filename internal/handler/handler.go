package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/handler/middleware"
)

func InitializationHandlers(ctx context.Context, cfg config.Config) (*gin.Engine, error) {
	// Создаем UrlHandler
	urlHandler, err := NewUrlHandler(ctx, cfg)
	if err != nil {
		return nil, err
	}

	gin.SetMode(cfg.App.GinMode)

	// Создаем новый роутер Gin
	router := gin.New()

	// Применяем middleware
	router.Use(middleware.GinMiddleware())

	// Регистрируем маршруты из UrlHandler
	urlHandler.RegisterRoutes(router)

	return router, nil
}
