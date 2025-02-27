package main

import (
	"context"
	"fmt"
	"leenwood/yandex-http/config"
	handlers "leenwood/yandex-http/internal/handler"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	ctx := context.Background()

	h, err := handlers.InitializationHandlers(ctx, cfg)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("0.0.0.0:%s", cfg.App.Port)
	err = http.ListenAndServe(url, h)
	if err != nil {
		panic(err)
	}
}
