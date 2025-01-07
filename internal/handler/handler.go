package handlers

import (
	"fmt"
	"leenwood/yandex-http/internal/handler/middleware"
	"net/http"
)

func InilizationHandlers() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.Handle(`/`, middleware.Middleware(http.HandlerFunc(mainPage)))
	return mux, nil
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Method: %s\r\n", req.Method)
	body += "Header =========================== \r\n"
	for k, v := range req.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "\r\n"
	body += "Query Params ===================== \r\n"
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}
	for k, v := range req.Form {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	res.Write([]byte(body))
}
