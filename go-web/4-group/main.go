package main

import (
  "net/http"

  "go-study/go-web/4-group/gee"
)

func main() {
  r := gee.New()
  r.GET("/index", func(c *gee.Context) {
    c.HTML(http.StatusOK, "<h1>Index Page</h1>")
  })
  v1 := r.Group("/v1")
  {
    v1.GET("/", func(c *gee.Context) {
      c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
    })

    v1.GET("/hello", func(c *gee.Context) {
      c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
    })
  }
  v2 := r.Group("/v2")
  {
    v2.GET("/hello/:name", func(c *gee.Context) {
      c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
    })

    v2.POST("/login", func(c *gee.Context) {
      c.JSON(http.StatusOK, gee.H{
        "username": c.PostForm("username"),
        "password": c.PostForm("password"),
      })
    })
  }

  r.Run(":9999")
}

/*
$ go run go-web/4-group/main.go
2023/08/22 11:13:28 Route  GET - /index
2023/08/22 11:13:28 Route  GET - /v1/
2023/08/22 11:13:28 Route  GET - /v1/hello
2023/08/22 11:13:28 Route  GET - /v2/hello/:name
2023/08/22 11:13:28 Route POST - /v2/login
$ curl localhost:9999/v1/hello?name=lb
hello lb, you're at /v1/hello
$ curl localhost:9999/v2/hello/lb
hello lb, you're at /v2/hello/lb
*/
