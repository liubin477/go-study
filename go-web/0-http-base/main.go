package main

import (
  "fmt"
  "log"
  "net/http"
)

func main() {
  // 两个路由绑定两个 Handler
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/hello", helloHandler)
  // 启动 Web 服务
  log.Fatal(http.ListenAndServe(":9999", nil))
}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
  for k, v := range req.Header {
    fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
  }
}

/*
$ go run go-web/0-http-base/main.go
$ curl localhost:9999
URL.Path = "/"
$ curl localhost:9999/hello
Header["User-Agent"] = ["curl/8.0.1"]
Header["Accept"] = ["*\/*"]
*/
