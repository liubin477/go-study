package main

import (
  "fmt"
  "net/http"

  "go-study/go-web/1-base/gee"
)

func main() {
  r := gee.New()
  r.GET("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
  })

  r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
    for k, v := range req.Header {
      fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
    }
  })

  r.Run(":9999")
}

/*
$ go run go-web/1-base/main.go
$ curl localhost:9999
URL.Path = "/"
$ curl localhost:9999/hello
Header["User-Agent"] = ["curl/8.0.1"]
Header["Accept"] = ["*\/*"]
*/
