package main

import (
  "net/http"

  "go-study/go-web/3-router/gee"
)

func main() {
  r := gee.New()
  r.GET("/", func(c *gee.Context) {
    c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
  })

  r.GET("/hello", func(c *gee.Context) {
    // expect /hello?name=lb
    c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
  })

  r.GET("/hello/:name", func(c *gee.Context) {
    // expect /hello/lb
    c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
  })

  r.GET("/assets/*filepath", func(c *gee.Context) {
    c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
  })

  r.Run(":9999")
}

/*
$ go run go-web/3-router/main.go
$ curl localhost:9999/hello
hello , you're at /hello
$ curl localhost:9999/hello/lb
hello lb, you're at /hello/lb
$ curl localhost:9999/assets/css/main.css
{"filepath":"css/main.css"}
*/
