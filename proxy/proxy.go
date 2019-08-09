package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewProxy returns reverse proxy which maps path to backendURL
func NewProxy(path, backendURL string) (http.Handler, error) {
	target, err := url.ParseRequestURI(backendURL)
	if err != nil {
		return nil, err
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("proxy req.URL: %s, target: %s", r.URL.String(), target.String())
		proxy := httputil.NewSingleHostReverseProxy(target)

		r.URL.Path = strings.TrimPrefix(r.URL.Path, path)

		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.Host = target.Host

		// Update the headers to allow for SSL redirection
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

		log.Printf("req.URL before serving proxy: %s", r.URL.String())

		// this is non-blocking because it starts new goroutine
		proxy.ServeHTTP(w, r)
	})

	return handler, nil
}
