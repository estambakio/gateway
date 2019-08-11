package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy returns http.Handler which maps source path on proxy to backend URL
func NewProxy(from, to string) (http.Handler, error) {
	target, err := url.ParseRequestURI(to)
	if err != nil {
		return nil, err
	}

	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Update the headers to allow for SSL redirection
		r.Host = target.Host
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

		// remove "from" value from r.URL.Path and proxy request through target URI
		http.StripPrefix(from, httputil.NewSingleHostReverseProxy(target)).ServeHTTP(w, r)
	})

	return proxyHandler, nil
}
