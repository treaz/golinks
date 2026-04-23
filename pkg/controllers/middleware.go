package controllers

import (
	"log/slog"
	"net/http"
	"net/url"
)

func (ac *AppController) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(
			"request",
			slog.String("remote", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}

func (ac *AppController) DecodeSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if decoded, err := url.PathUnescape(r.URL.EscapedPath()); err == nil {
			r.URL.Path = decoded
			r.URL.RawPath = ""
		}
		next.ServeHTTP(w, r)
	})
}
