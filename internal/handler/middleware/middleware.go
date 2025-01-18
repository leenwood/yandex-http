package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
		fmt.Printf("Request url - %s, Request time - %s\r\n", r.RequestURI, formatted)
		next.ServeHTTP(w, r)
	})
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Захватываем текущее время
		t := time.Now()
		formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		// Логируем информацию о запросе
		fmt.Printf("Request url - %s, Request time - %s\r\n", c.Request.RequestURI, formatted)

		// Продолжаем обработку следующего middleware или хэндлера
		c.Next()
	}
}
