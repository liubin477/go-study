package main

import (
  "fmt"
  "log"
  "net/http"
)

// Engine is the uni handler for all requests
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  switch req.URL.Path {
  case "/":
    fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
  case "/hello":
    for k, v := range req.Header {
      fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
    }
  default:
    fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
  }
}

func main() {
  engine := new(Engine)
  log.Fatal(http.ListenAndServe(":9999", engine))
}

/*
$ go run go-web\0-http-base\main-handler.go
$ curl localhost:9999
URL.Path = "/"
$ curl localhost:9999/hello
Header["User-Agent"] = ["curl/8.0.1"]
Header["Accept"] = ["*\/*"]
*/
