package main

import (
	handlers "leenwood/yandex-http/internal/handler"
	"net/http"
)

func main() {
	h, err := handlers.InilizationHandlers()
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(`:8000`, h)
	if err != nil {
		panic(err)
	}
}
