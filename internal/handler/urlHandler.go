package handlers

import (
	"context"
	"fmt"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/usecase"
	"leenwood/yandex-http/internal/usecase/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UrlHandler struct {
	us usecase.UrlUseCaseInterface
}

func NewUrlHandler(ctx context.Context, cfg config.Config) (*UrlHandler, error) {
	us, err := usecase.NewUrlUseCase(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &UrlHandler{us: us}, nil
}

func (uh *UrlHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/", uh.CreateShortUrl)
	router.GET("/healthz", uh.CheckHealthz)
}
func (uh *UrlHandler) CreateShortUrl(c *gin.Context) {
	var req dto.CreateShortUrlRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data: " + err.Error()})
		return
	}

	var (
		response interface{}
		err      error
	)

	// Обработка в зависимости от наличия ID
	if req.Id != "" {
		response, err = uh.handleCustomIdRequest(req)
	} else {
		response, err = uh.handleDefaultRequest(req)
	}

	// Проверка на ошибки
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create short URL: " + err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, response)
}

// Обработка запроса с пользовательским ID
func (uh *UrlHandler) handleCustomIdRequest(req dto.CreateShortUrlRequest) (interface{}, error) {
	request := dto.CreateShortUrlWithCustomIdRequest{
		Url: req.Url,
		Id:  req.Id,
	}
	return uh.us.CreateShortUrlWithCustomId(request)
}

// Обработка запроса без пользовательского ID
func (uh *UrlHandler) handleDefaultRequest(req dto.CreateShortUrlRequest) (interface{}, error) {
	request := dto.CreateShortUrlUseCaseRequest{
		Url: req.Url,
	}
	return uh.us.CreateShortUrl(request)
}

func (uh *UrlHandler) CheckHealthz(c *gin.Context) {
	body := fmt.Sprintf("Method: %s\r\n", c.Request.Method)
	body += "Header =========================== \r\n"
	for k, v := range c.Request.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "\r\n"
	body += "Query Params ===================== \r\n"
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse query params: " + err.Error()})
		return
	}
	for k, v := range c.Request.Form {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}

	c.String(http.StatusOK, body)
}
