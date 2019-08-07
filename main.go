package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// web url: https://domain.com/project/app/path?query
// ingress sends to proxy: http://proxy/app/path?query
// proxy should pass request to http://app/path?query
var urls = map[string]string{
	"http://localhost:3000/app": "http://google.com",
}

func serveReverseProxy(rw http.ResponseWriter, req *http.Request) {
	var match string
	for k := range urls {
		kURL, err := url.Parse(k)
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(req.URL.Path, kURL.Path) {
			match = k
			break
		}
	}

	// dummy response if nothing matched
	if match == "" {
		http.Error(rw, "No match in proxy rules", http.StatusBadRequest)
		return
	}

	v := urls[match]
	target, err := url.Parse(v)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Target: %s\n", target.String())

	proxy := httputil.NewSingleHostReverseProxy(target)

	req.URL.Host = target.Host
	req.URL.Scheme = target.Scheme

	source, err := url.Parse(match)
	if err != nil {
		panic(err)
	}
	req.URL.Path = strings.Join([]string{target.Path, strings.TrimPrefix(req.URL.Path, source.Path)}, "/")

	// Update the headers to allow for SSL redirection
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = target.Host

	fmt.Printf("Req to proxy URL: %s\n", req.URL.String())

	proxy.ServeHTTP(rw, req)
}

func main() {
	http.HandleFunc("/", serveReverseProxy)
	log.Fatal(http.ListenAndServe(":"+port(), nil))
}

func port() string {
	if v, ok := os.LookupEnv("PORT"); ok {
		return v
	}
	return "3000"
}
