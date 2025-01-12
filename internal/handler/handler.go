package handlers

import (
	"context"
	"fmt"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/handler/middleware"
	"net/http"
)

func InitializationHandlers(ctx context.Context, cfg config.Config) (*http.ServeMux, error) {
	url, err := NewUrlHandler(ctx, cfg)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle(`/`, middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		url.CreateShortUrl(w, r)
	})))
	mux.Handle(`/healthz`, middleware.Middleware(http.HandlerFunc(checkHealthz)))
	return mux, nil
}

func checkHealthz(res http.ResponseWriter, req *http.Request) {
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
	_, err = res.Write([]byte(body))
	if err != nil {
		return
	}
}
