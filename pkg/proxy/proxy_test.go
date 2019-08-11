package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// /app http://app/app-context-path
// - req to proxy: /app/something
// - backend req: http://app/app-context-path/something
// /app2 http://app2
// - req to proxy: /app2/something
// - backend req: http://app2/something
func Test_NewProxy(t *testing.T) {
	// backend server responds with its request URI to indicate what request it has received
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.RequestURI())
	}))
	defer backendServer.Close()

	tests := []struct {
		// if proxy is deployed as http://proxy.com
		// and proxyPath is /myapp, backendBaseURI is /some_backend_path
		// then it means that all requests to http://proxy.com/myapp/some_page.html are routed to
		// http://backend_domain.com/some_backend_path/some_page.html
		// backend_domain.com here is known in tests because it's set up earlier;
		// function NewProxy receives full URL to backend, for example in our case:
		// proxy, err := NewProxy("/myapp", "http://backend_domain.com/some_backend_path")
		proxyPath             string
		backendBaseURI        string
		requestURI            string
		expectedBackendReqURI string
	}{
		{"/app", "/context_path", "/app/payload", "/context_path/payload"},
		{"/app", "/context_path", "/app/payload?a=b&c=45", "/context_path/payload?a=b&c=45"},
		{"/app", "/", "/app/payload?a=b&c=45", "/payload?a=b&c=45"},
		{"/myapp", "/some_path", "/myapp/info.html", "/some_path/info.html"},
	}

	for _, test := range tests {
		t.Run(test.proxyPath, func(t *testing.T) {
			proxy, err := NewProxy(test.proxyPath, backendServer.URL+test.backendBaseURI)
			if err != nil {
				t.Errorf("failed to create test proxy: %v", err)
			}

			// test both http and https proxy
			frontends := []struct {
				schema string
				server *httptest.Server
			}{
				{schema: "http", server: httptest.NewServer(proxy)},
				{schema: "https", server: httptest.NewTLSServer(proxy)},
			}

			for _, frontend := range frontends {
				t.Run(frontend.schema, func(t *testing.T) {
					frontendProxy := frontend.server
					defer frontendProxy.Close()

					// this is required for TLS test server to not fail on certificates
					// for http server http.Get will do, but this approach handles both test cases
					client := frontendProxy.Client()

					// TODO test also other request methods (POST etc.)
					resp, err := client.Get(frontendProxy.URL + test.requestURI)
					if err != nil {
						t.Errorf("failed to get %s: %v", frontendProxy.URL+test.requestURI, err)
					}
					defer resp.Body.Close()

					b, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						t.Error(err)
					}

					result := string(b)

					if result != test.expectedBackendReqURI {
						t.Errorf("b: %s, but expected %s", result, test.expectedBackendReqURI)
					}
				})
			}

		})
	}

	// should fail if passed invalid url as backend
	_, err := NewProxy("/app", "1234")
	if err == nil {
		t.Errorf("did not fail for invalid backend url")
	}
}
