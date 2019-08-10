package main

import (
	"log"
	"net/http"

	"github.com/estambakio/gateway/pkg/proxy"
)

func main() {
	mapping := map[string]string{
		// how to run example:
		// go run main.go
		// curl -v http://localhost:3000/task/1 - see JSON with id = 1
		"/task/": "https://jsonplaceholder.typicode.com/todos/",
	}

	for from, to := range mapping {
		handler, err := proxy.NewProxy(from, to)
		if err != nil {
			panic(err)
		}
		// register proxy handler for every path
		http.Handle(from, handler)
	}

	log.Fatal(http.ListenAndServe(":3000", nil))
}
